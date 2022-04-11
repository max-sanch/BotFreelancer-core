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

	query := fmt.Sprintf("SELECT * FROM %s WHERE tg_id = %v", usersTable, tgID)
	row := r.db.QueryRow(query)
	if err := row.Scan(&user.ID, &user.TGID, &user.Username); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE user_id = %v", userSettingsTable, user.ID)
	row = r.db.QueryRow(query)
	if err := row.Scan(&setting.ID, &setting.IsSafeDeal, &setting.IsBudget, &setting.IsTerm); err != nil {
		return core.UserResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE user_setting_id = %v", userCategoriesTable, setting.ID)
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
	var id int

	query := fmt.Sprintf(`
	BEGIN;
	INSERT INTO %s (tg_id, username) VALUES
	(%v, '%s') RETURNING id;
	
	WITH user_i AS (
	  SELECT id FROM %s WHERE tg_id = %v
	), user_setting AS (
	  INSERT INTO %s (user_id, is_safe_deal, is_budget, is_term) VALUES
	  ((SELECT id FROM user_i), %t, %t, %t) RETURNING id
	)
	
	SELECT insert_user_categories(ARRAY%v, (SELECT id FROM user_setting));
	COMMIT;
	`, usersTable, userInput.TGID, userInput.Username, usersTable, userInput.TGID,
		userSettingsTable, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm,
		fmt.Sprintf("[%s]", arrayToString(userInput.Setting.Categories, ", ")))

	row := r.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserPostgres) UpdateUser(userInput core.UserInput) (int, error) {
	var id int

	query := fmt.Sprintf(`
	BEGIN;
	UPDATE %s SET username = '%s' WHERE tg_id = %v RETURNING id;
	
	WITH user_i AS (
	  SELECT id FROM %s WHERE tg_id = %v
	), user_setting_d AS (
	  UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t
	  WHERE user_id = (SELECT id FROM user_i) RETURNING id
	), user_setting_i AS (
	  DELETE FROM %s WHERE user_setting_id = (SELECT id FROM user_setting_d)
	  RETURNING user_setting_id AS id
	)
	
	SELECT insert_user_categories(ARRAY%v, (SELECT id FROM user_setting_i LIMIT 1));
	COMMIT;
	`, usersTable, userInput.Username, userInput.TGID, usersTable, userInput.TGID,
		userSettingsTable, *userInput.Setting.IsSafeDeal, *userInput.Setting.IsBudget, *userInput.Setting.IsTerm,
		userCategoriesTable ,fmt.Sprintf("[%s]", arrayToString(userInput.Setting.Categories, ", ")))

	row := r.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
