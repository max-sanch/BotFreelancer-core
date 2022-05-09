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
	var channel core.ChannelResponse
	var settingId int

	query := fmt.Sprintf("SELECT * FROM %s WHERE api_id = $1", channelsTable)
	if err := r.db.Get(&channel, query, apiId); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE channel_id = $1",
		channelSettingsTable)
	row := r.db.QueryRow(query, channel.Id)
	if err := row.Scan(&settingId, &channel.Setting.IsSafeDeal, &channel.Setting.IsBudget, &channel.Setting.IsTerm); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE channel_setting_id = $1", channelCategoriesTable)
	if err := r.db.Select(&channel.Setting.Categories, query, settingId); err != nil {
		return core.ChannelResponse{}, err
	}

	return channel, nil
}

func (r *ChannelPostgres) Create(channelInput core.ChannelInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var channelId, channelSettingId int
	createChannelQuery := fmt.Sprintf("INSERT INTO %s (api_id, api_hash, name) VALUES ($1, $2, $3) RETURNING id;",
		channelsTable)

	row := tx.QueryRow(createChannelQuery, channelInput.ApiId, channelInput.ApiHash, channelInput.Name)
	if err := row.Scan(&channelId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	createChannelSettingQuery := fmt.Sprintf("INSERT INTO %s (channel_id, is_safe_deal, is_budget, is_term) VALUES ($1, $2, $3, $4) RETURNING id;",
		channelSettingsTable)

	if channelInput.Setting.IsSafeDeal == nil || channelInput.Setting.IsBudget == nil || channelInput.Setting.IsTerm == nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
	}

	row = tx.QueryRow(createChannelSettingQuery, channelId, *channelInput.Setting.IsSafeDeal,
		*channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm)
	if err := row.Scan(&channelSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES ($1, $2);",
			channelCategoriesTable)

		if _, err := tx.Exec(createChannelCategoryQuery, channelSettingId, categoryId); err != nil {
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
	updateChannelQuery := fmt.Sprintf("UPDATE %s SET api_hash = $1, name = $2 WHERE api_id = $3 RETURNING id;",
		channelsTable)

	row := tx.QueryRow(updateChannelQuery, channelInput.ApiHash, channelInput.Name, channelInput.ApiId)
	if err := row.Scan(&channelId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	updateChannelSettingQuery := fmt.Sprintf("UPDATE %s SET is_safe_deal = $1, is_budget = $2, is_term = $3 WHERE channel_id = $4 RETURNING id;",
		channelSettingsTable)

	if channelInput.Setting.IsSafeDeal == nil || channelInput.Setting.IsBudget == nil || channelInput.Setting.IsTerm == nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
	}

	row = tx.QueryRow(updateChannelSettingQuery, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm, channelId)
	if err := row.Scan(&channelSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	deleteChannelCategoryQuery := fmt.Sprintf("DELETE FROM %s WHERE channel_setting_id = $1;", channelCategoriesTable)

	if _, err := tx.Exec(deleteChannelCategoryQuery, channelSettingId); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	for _, categoryId := range channelInput.Setting.Categories {
		createChannelCategoryQuery := fmt.Sprintf("INSERT INTO %s (channel_setting_id, category_id) VALUES ($1, $2);",
			channelCategoriesTable)

		if _, err := tx.Exec(createChannelCategoryQuery, channelSettingId, categoryId); err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
	}

	return channelId, tx.Commit()
}

func (r *ChannelPostgres) Delete(apiId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE api_id = $1;", channelsTable)

	if _, err := r.db.Exec(query, apiId); err != nil {
		return err
	}
	return nil
}
