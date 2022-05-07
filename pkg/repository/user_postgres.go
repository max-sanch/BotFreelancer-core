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

	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_id = %d", usersTable, tgId)
	if err := r.db.Get(&user, query); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE user_id = %d",
		userSettingsTable, user.Id)
	row := r.db.QueryRow(query)
	if err := row.Scan(&settingId, &user.Setting.IsSafeDeal, &user.Setting.IsBudget, &user.Setting.IsTerm); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE user_setting_id = %d", userCategoriesTable, settingId)
	if err := r.db.Select(&user.Setting.Categories, query); err != nil {
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
	createUserQuery := fmt.Sprintf("INSERT INTO %s (tg_id, username) VALUES (%d, '%s') RETURNING id;",
		usersTable, userInput.TgId, userInput.Username)

	row := tx.QueryRow(createUserQuery)
	if err := row.Scan(&userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createUserSettingQuery := fmt.Sprintf("INSERT INTO %s (user_id, is_safe_deal, is_budget, is_term) VALUES (%d, %t, %t, %t) RETURNING id;",
		userSettingsTable, userId, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm)

	row = tx.QueryRow(createUserSettingQuery)
	if err := row.Scan(&userSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES (%d, %d);",
			userCategoriesTable, userSettingId, categoryId)

		if _, err := tx.Exec(createUserCategoryQuery); err != nil {
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
	updateUserQuery := fmt.Sprintf("UPDATE %s SET username = '%s' WHERE tg_id = %d RETURNING id;",
		usersTable, userInput.Username, userInput.TgId)

	row := tx.QueryRow(updateUserQuery)
	if err := row.Scan(&userId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	updateUserSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t WHERE user_id = %d RETURNING id;",
		userSettingsTable, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm, userId)

	row = tx.QueryRow(updateUserSettingQuery)
	if err := row.Scan(&userSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteUserCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE user_setting_id = %d;",
		userCategoriesTable, userSettingId)

	if _, err := tx.Exec(deleteUserCategoryQuery); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES (%d, %d);",
			userCategoriesTable, userSettingId, categoryId)

		if _, err := tx.Exec(createUserCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}
	return userId, tx.Commit()
}
