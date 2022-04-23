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

func (r *ChannelPostgres) GetChannel(apiID int) (core.ChannelResponse, error) {
	var channel core.Channel
	var setting core.Setting
	var category int
	var categories []int

	query := fmt.Sprintf("SELECT * FROM %s WHERE api_id = %d", channelsTable, apiID)
	row := r.db.QueryRow(query)
	if err := row.Scan(&channel.ID, &channel.APIID, &channel.APIHash, &channel.Name); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE channel_id = %d",
		channelSettingsTable, channel.ID)
	row = r.db.QueryRow(query)
	if err := row.Scan(&setting.ID, &setting.IsSafeDeal, &setting.IsBudget, &setting.IsTerm); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE channel_setting_id = %d",
		channelCategoriesTable, setting.ID)
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
		ID: channel.ID,
		APIID: channel.APIID,
		APIHash: channel.APIHash,
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

func(r *ChannelPostgres) CreateChannel(channelInput core.ChannelInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var channelID, channelSettingID int
	createChannelQuery := fmt.Sprintf("INSERT INTO %s (api_id, api_hash, name) VALUES (%d, '%s', '%s') RETURNING id;",
		channelsTable, channelInput.APIID, channelInput.APIHash, channelInput.Name)

	row := tx.QueryRow(createChannelQuery)
	if err := row.Scan(&channelID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createChannelSettingQuery := fmt.Sprintf("INSERT INTO %s (channel_id, is_safe_deal, is_budget, is_term) VALUES (%d, %t, %t, %t) RETURNING id;",
		channelSettingsTable, channelID, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm)

	row = tx.QueryRow(createChannelSettingQuery)
	if err := row.Scan(&channelSettingID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryID := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES (%d, %d);",
			channelCategoriesTable, channelSettingID, categoryID)

		if _, err := tx.Exec(createChannelCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return channelID, tx.Commit()
}

func (r *ChannelPostgres) UpdateChannel(channelInput core.ChannelInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var channelID, channelSettingID int
	updateChannelQuery := fmt.Sprintf("UPDATE %s SET api_hash = '%s', name = '%s' WHERE api_id = %d RETURNING id;",
		channelsTable, channelInput.APIHash, channelInput.Name, channelInput.APIID)

	row := tx.QueryRow(updateChannelQuery)
	if err := row.Scan(&channelID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	updateChannelSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t WHERE channel_id = %d RETURNING id;",
		channelSettingsTable, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm, channelID)

	row = tx.QueryRow(updateChannelSettingQuery)
	if err := row.Scan(&channelSettingID); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteChannelCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE channel_setting_id = %d;",
		channelCategoriesTable, channelSettingID)

	if _, err := tx.Exec(deleteChannelCategoryQuery); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryID := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES (%d, %d);",
			channelCategoriesTable, channelSettingID, categoryID)

		if _, err := tx.Exec(createChannelCategoryQuery); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return channelID, tx.Commit()
}

func (r *ChannelPostgres) DeleteChannel(apiID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE api_id = %d;", channelsTable, apiID)

	row := r.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
