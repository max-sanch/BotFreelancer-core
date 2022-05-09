package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	core "github.com/max-sanch/BotFreelancer-core"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetByTgId(tgId int) (core.UserResponse, error) {
	var user core.UserResponse
	var settingId int

	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_id = $1", usersTable)
	if err := r.db.Get(&user, query, tgId); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE user_id = $1",
		userSettingsTable)
	row := r.db.QueryRow(query, user.Id)
	if err := row.Scan(&settingId, &user.Setting.IsSafeDeal, &user.Setting.IsBudget, &user.Setting.IsTerm); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE user_setting_id = $1", userCategoriesTable)
	if err := r.db.Select(&user.Setting.Categories, query, settingId); err != nil {
		return core.UserResponse{}, err
	}

	return user, nil
}

func (r *UserPostgres) Create(userInput core.UserInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userId, userSettingId int
	createUserQuery := fmt.Sprintf("INSERT INTO %s (tg_id, username) VALUES ($1, $2) RETURNING id;", usersTable)

	row := tx.QueryRow(createUserQuery, userInput.TgId, userInput.Username)
	if err := row.Scan(&userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	if userInput.Setting.IsSafeDeal == nil || userInput.Setting.IsBudget == nil || userInput.Setting.IsTerm == nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
	}

	createUserSettingQuery := fmt.Sprintf("INSERT INTO %s (user_id, is_safe_deal, is_budget, is_term) VALUES ($1, $2, $3, $4) RETURNING id;",
		userSettingsTable)

	row = tx.QueryRow(createUserSettingQuery, userId, *userInput.Setting.IsSafeDeal,
		*userInput.Setting.IsBudget, *userInput.Setting.IsTerm)
	if err := row.Scan(&userSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES ($1, $2);",
			userCategoriesTable)

		if _, err := tx.Exec(createUserCategoryQuery, userSettingId, categoryId); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return userId, tx.Commit()
}

func (r *UserPostgres) Update(userInput core.UserInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userId, userSettingId int
	updateUserQuery := fmt.Sprintf("UPDATE %s SET username = $1 WHERE tg_id = $2 RETURNING id;", usersTable)

	row := tx.QueryRow(updateUserQuery, userInput.Username, userInput.TgId)
	if err := row.Scan(&userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	if userInput.Setting.IsSafeDeal == nil || userInput.Setting.IsBudget == nil || userInput.Setting.IsTerm == nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
	}

	updateUserSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = $1, is_budget = $2, is_term = $3 WHERE user_id = $4 RETURNING id;",
		userSettingsTable)

	row = tx.QueryRow(updateUserSettingQuery, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm, userId)
	if err := row.Scan(&userSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteUserCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE user_setting_id = $1;", userCategoriesTable)

	if _, err := tx.Exec(deleteUserCategoryQuery, userSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES ($1, $2);",
			userCategoriesTable)

		if _, err := tx.Exec(createUserCategoryQuery, userSettingId, categoryId); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}
	return userId, tx.Commit()
}
