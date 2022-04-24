package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	core "github.com/max-sanch/BotFreelancer-core"
)

type ChannelPostgres struct {
	db *sqlx.DB
}

func NewChannelPostgres(db *sqlx.DB) *ChannelPostgres {
	return &ChannelPostgres{db: db}
}

func (r *ChannelPostgres) GetByApiId(apiId int) (core.ChannelResponse, error) {
	var channel ChannelObject
	var setting SettingObject
	var category int
	var categories []int

	query := fmt.Sprintf("SELECT * FROM %s WHERE api_id = %d", channelsTable, apiId)
	row := r.db.QueryRow(query)
	if err := row.Scan(&channel.Id, &channel.ApiId, &channel.ApiHash, &channel.Name); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE channel_id = %d",
		channelSettingsTable, channel.Id)
	row = r.db.QueryRow(query)
	if err := row.Scan(&setting.Id, &setting.IsSafeDeal, &setting.IsBudget, &setting.IsTerm); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE channel_setting_id = %d",
		channelCategoriesTable, setting.Id)
	rows, err := r.db.Query(query)
	if err != nil {
		return core.ChannelResponse{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&category); err != nil {
			return core.ChannelResponse{}, err
		}

		categories = append(categories, category)
	}

	channelResponse := core.ChannelResponse{
		Id: channel.Id,
		ApiId: channel.ApiId,
		ApiHash: channel.ApiHash,
		Name: channel.Name,
		Setting: core.SettingResponse{
			IsSafeDeal: setting.IsSafeDeal,
			IsBudget: setting.IsBudget,
			IsTerm: setting.IsTerm,
			Categories: categories,
		},
	}

	return channelResponse, nil
}

func(r *ChannelPostgres) Create(channelInput core.ChannelInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var channelId, channelSettingId int
	createChannelQuery := fmt.Sprintf("INSERT INTO %s (api_id, api_hash, name) VALUES (%d, '%s', '%s') RETURNING id;",
		channelsTable, channelInput.ApiId, channelInput.ApiHash, channelInput.Name)

	row := tx.QueryRow(createChannelQuery)
	if err := row.Scan(&channelId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createChannelSettingQuery := fmt.Sprintf("INSERT INTO %s (channel_id, is_safe_deal, is_budget, is_term) VALUES (%d, %t, %t, %t) RETURNING id;",
		channelSettingsTable, channelId, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm)

	row = tx.QueryRow(createChannelSettingQuery)
	if err := row.Scan(&channelSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES (%d, %d);",
			channelCategoriesTable, channelSettingId, categoryId)

		if _, err := tx.Exec(createChannelCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return channelId, tx.Commit()
}

func (r *ChannelPostgres) Update(channelInput core.ChannelInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var channelId, channelSettingId int
	updateChannelQuery := fmt.Sprintf("UPDATE %s SET api_hash = '%s', name = '%s' WHERE api_id = %d RETURNING id;",
		channelsTable, channelInput.ApiHash, channelInput.Name, channelInput.ApiId)

	row := tx.QueryRow(updateChannelQuery)
	if err := row.Scan(&channelId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	updateChannelSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t WHERE channel_id = %d RETURNING id;",
		channelSettingsTable, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm, channelId)

	row = tx.QueryRow(updateChannelSettingQuery)
	if err := row.Scan(&channelSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteChannelCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE channel_setting_id = %d;",
		channelCategoriesTable, channelSettingId)

	if _, err := tx.Exec(deleteChannelCategoryQuery); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES (%d, %d);",
			channelCategoriesTable, channelSettingId, categoryId)

		if _, err := tx.Exec(createChannelCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return channelId, tx.Commit()
}

func (r *ChannelPostgres) Delete(apiId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE api_id = %d;", channelsTable, apiId)

	row := r.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
