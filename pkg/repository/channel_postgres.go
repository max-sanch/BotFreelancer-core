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

	query := fmt.Sprintf("SELECT * FROM %s WHERE api_id = %v", channelsTable, apiID)
	row := r.db.QueryRow(query)
	if err := row.Scan(&channel.ID, &channel.APIID, &channel.APIHash, &channel.Name); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT id, is_safe_deal, is_budget, is_term FROM %s WHERE channel_id = %v", channelSettingsTable, channel.ID)
	row = r.db.QueryRow(query)
	if err := row.Scan(&setting.ID, &setting.IsSafeDeal, &setting.IsBudget, &setting.IsTerm); err != nil {
		return core.ChannelResponse{}, err
	}

	query = fmt.Sprintf("SELECT category_id FROM %s WHERE channel_setting_id = %v", channelCategoriesTable, setting.ID)
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

func (r *ChannelPostgres) CreateChannel(channelInput core.ChannelInput) (int, error) {
	var id int

	query := fmt.Sprintf(`
	BEGIN;
	INSERT INTO %s (api_id, api_hash, name) VALUES
	(%v, '%s', '%s') RETURNING id;
	
	WITH channel AS (
	  SELECT id FROM %s WHERE api_id = %v
	), channel_setting AS (
	  INSERT INTO %s (channel_id, is_safe_deal, is_budget, is_term) VALUES
	  ((SELECT id FROM channel), %t, %t, %t) RETURNING id
	)
	
	SELECT insert_channel_categories(ARRAY%v, (SELECT id FROM channel_setting));
	COMMIT;
	`, channelsTable, channelInput.APIID, channelInput.APIHash, channelInput.Name, channelsTable, channelInput.APIID,
	channelSettingsTable, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm,
	fmt.Sprintf("[%s]", arrayToString(channelInput.Setting.Categories, ", ")))

	row := r.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ChannelPostgres) UpdateChannel(channelInput core.ChannelInput) (int, error) {
	var id int

	query := fmt.Sprintf(`
	BEGIN;
	UPDATE %s SET api_hash = '%s', name = '%s' WHERE api_id = %v RETURNING id;
	
	WITH channel AS (
	  SELECT id FROM %s WHERE api_id = %v
	), channel_setting_d AS (
	  UPDATE %s SET is_safe_deal = %t, is_budget = %t, is_term = %t
	  WHERE channel_id = (SELECT id FROM channel) RETURNING id
	), channel_setting_i AS (
	  DELETE FROM %s WHERE channel_setting_id = (SELECT id FROM channel_setting_d)
	  RETURNING channel_setting_id AS id
	)
	
	SELECT insert_channel_categories(ARRAY%v, (SELECT id FROM channel_setting_i LIMIT 1));
	COMMIT;
	`, channelsTable, channelInput.APIHash, channelInput.Name, channelInput.APIID, channelsTable, channelInput.APIID,
		channelSettingsTable, *channelInput.Setting.IsSafeDeal, *channelInput.Setting.IsBudget, *channelInput.Setting.IsTerm,
		channelCategoriesTable ,fmt.Sprintf("[%s]", arrayToString(channelInput.Setting.Categories, ", ")))

	row := r.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ChannelPostgres) DeleteChannel(apiID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE api_id = %v;", channelsTable, apiID)

	row := r.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
