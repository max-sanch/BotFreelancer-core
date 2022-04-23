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

func (r *UserPostgres) GetUser(tgID int) (core.UserResponse, error) {
	var user core.User
	var setting core.Setting
	var category int
	var categories []int

	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_id = %d", usersTable, tgID)
	row := r.db.QueryRow(query)
	if err := row.Scan(&user.ID, &user.TGID, &user.Username); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE user_id = %d", userSettingsTable, user.ID)
	row = r.db.QueryRow(query)
	if err := row.Scan(&setting.ID, &setting.IsSafeDeal, &setting.IsBudget, &setting.IsTerm); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE user_setting_id = %d", userCategoriesTable, setting.ID)
	rows, err := r.db.Query(query)
	if err != nil {
		return core.UserResponse{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&category); err != nil {
			return core.UserResponse{}, err
		}

		categories = append(categories, category)
	}

	userResponse := core.UserResponse{
		ID: user.ID,
		TGID: user.TGID,
		Username: user.Username,
		Setting: core.SettingResponse{
			IsSafeDeal: setting.IsSafeDeal,
			IsBudget: setting.IsBudget,
			IsTerm: setting.IsTerm,
			Categories: categories,
		},
	}

	return userResponse, nil
}

func (r *UserPostgres) CreateUser(userInput core.UserInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userID, userSettingID int
	createUserQuery := fmt.Sprintf("INSERT INTO %s (tg_id, username) VALUES (%d, '%s') RETURNING id;",
		usersTable, userInput.TGID, userInput.Username)

	row := tx.QueryRow(createUserQuery)
	if err := row.Scan(&userID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createUserSettingQuery := fmt.Sprintf("INSERT INTO %s (user_id, is_safe_deal, is_budget, is_term) VALUES (%d, %t, %t, %t) RETURNING id;",
		userSettingsTable, userID, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm)

	row = tx.QueryRow(createUserSettingQuery)
	if err := row.Scan(&userSettingID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryID := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES (%d, %d);",
			userCategoriesTable, userSettingID, categoryID)

		if _, err := tx.Exec(createUserCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return userID, tx.Commit()
}

func (r *UserPostgres) UpdateUser(userInput core.UserInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var userID, userSettingID int
	updateUserQuery := fmt.Sprintf("UPDATE %s SET username = '%s' WHERE tg_id = %d RETURNING id;",
		usersTable, userInput.Username, userInput.TGID)

	row := tx.QueryRow(updateUserQuery)
	if err := row.Scan(&userID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	updateUserSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t WHERE user_id = %d RETURNING id;",
		userSettingsTable, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm, userID)

	row = tx.QueryRow(updateUserSettingQuery)
	if err := row.Scan(&userSettingID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteUserCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE user_setting_id = %d;",
		userCategoriesTable, userSettingID)

	if _, err := tx.Exec(deleteUserCategoryQuery); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryID := range userInput.Setting.Categories {
		createUserCategoryQuery := fmt.Sprintf("INSERT INTO %s (user_setting_id, category_id) VALUES (%d, %d);",
			userCategoriesTable, userSettingID, categoryID)

		if _, err := tx.Exec(createUserCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}
	return userID, tx.Commit()
}
