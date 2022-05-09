package repository

import (
	"database/sql"
	"errors"
	"testing"

	core "github.com/max-sanch/BotFreelancer-core"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestChannelPostgres_GetByApiId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChannelPostgres(db)

	type args struct {
		apiId int
	}

	type mockBehavior func(args args)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         core.ChannelResponse
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				apiId: 1111,
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "api_id", "api_hash", "name"}).
					AddRow(1, 1111, "hash1111", "channel-1")

				mock.ExpectQuery("SELECT (.+) FROM channels WHERE (.+)").
					WithArgs(args.apiId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id", "is_safe_deal", "is_budget", "is_term"}).
					AddRow(1, true, true, true)

				mock.ExpectQuery("SELECT (.+) FROM channel_settings WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"category_id"}).AddRow(1).AddRow(2)

				mock.ExpectQuery("SELECT (.+) FROM channel_categories WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			want: core.ChannelResponse{
				Id:      1,
				ApiId:   1111,
				ApiHash: "hash1111",
				Name:    "channel-1",
				Setting: core.SettingResponse{
					IsSafeDeal: true,
					IsBudget:   true,
					IsTerm:     true,
					Categories: []int{1, 2},
				},
			},
		},
		{
			name: "Not Found",
			args: args{
				apiId: 1111,
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "api_id", "api_hash", "name"})

				mock.ExpectQuery("SELECT (.+) FROM channels WHERE (.+)").
					WithArgs(args.apiId).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTables {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetByApiId(testCase.args.apiId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChannelPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChannelPostgres(db)
	isFalse := false

	type args struct {
		channel core.ChannelInput
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO channels").WithArgs(
					args.channel.ApiId, args.channel.ApiHash, args.channel.Name).WillReturnRows(rows)

				channelSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(channelSettingId)
				mock.ExpectQuery("INSERT INTO channel_settings").WithArgs(
					id, args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm).WillReturnRows(rows)

				for _, categoryId := range args.channel.Setting.Categories {
					mock.ExpectExec("INSERT INTO channel_categories").WithArgs(
						channelSettingId, categoryId).WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO channels").WithArgs(
					args.channel.ApiId, args.channel.ApiHash, args.channel.Name).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: nil,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO channels").WithArgs(
					args.channel.ApiId, args.channel.ApiHash, args.channel.Name).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO channel_settings").WithArgs(
					id, args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "3nd Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO channels").WithArgs(
					args.channel.ApiId, args.channel.ApiHash, args.channel.Name).WillReturnRows(rows)

				channelSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(channelSettingId)
				mock.ExpectQuery("INSERT INTO channel_settings").WithArgs(
					id, args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO channel_categories").
					WithArgs(channelSettingId, args.channel.Setting.Categories[0]).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.channel)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestChannelPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChannelPostgres(db)
	isFalse := false

	type args struct {
		channel core.ChannelInput
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("UPDATE channels SET (.+) WHERE (.+)").WithArgs(
					args.channel.ApiHash, args.channel.Name, args.channel.ApiId).WillReturnRows(rows)

				channelSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(channelSettingId)
				mock.ExpectQuery("UPDATE channel_settings SET (.+) WHERE (.+)").WithArgs(
					args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM channel_categories WHERE (.+)").WithArgs(
					channelSettingId).WillReturnResult(sqlmock.NewResult(0, 1))

				for _, categoryId := range args.channel.Setting.Categories {
					mock.ExpectExec("INSERT INTO channel_categories").WithArgs(
						channelSettingId, categoryId).WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("UPDATE channels SET (.+) WHERE (.+)").WithArgs(
					args.channel.ApiHash, args.channel.Name, args.channel.ApiId).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   nil,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("UPDATE channels SET (.+) WHERE (.+)").WithArgs(
					args.channel.ApiHash, args.channel.Name, args.channel.ApiId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("UPDATE channel_settings SET (.+) WHERE (.+)").WithArgs(
					args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("UPDATE channels SET (.+) WHERE (.+)").WithArgs(
					args.channel.ApiHash, args.channel.Name, args.channel.ApiId).WillReturnRows(rows)

				channelSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(channelSettingId)
				mock.ExpectQuery("UPDATE channel_settings SET (.+) WHERE (.+)").WithArgs(
					args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM channel_categories WHERE (.+)").WithArgs(
					channelSettingId).WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "3nd Empty Fields",
			args: args{
				channel: core.ChannelInput{
					ApiId:   1111,
					ApiHash: "hash1111",
					Name:    "channel-1",
					Setting: core.SettingInput{
						IsSafeDeal: &isFalse,
						IsBudget:   &isFalse,
						IsTerm:     &isFalse,
						Categories: []int{1, 2},
					},
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("UPDATE channels SET (.+) WHERE (.+)").WithArgs(
					args.channel.ApiHash, args.channel.Name, args.channel.ApiId).WillReturnRows(rows)

				channelSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(channelSettingId)
				mock.ExpectQuery("UPDATE channel_settings SET (.+) WHERE (.+)").WithArgs(
					args.channel.Setting.IsSafeDeal, args.channel.Setting.IsBudget,
					args.channel.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM channel_categories WHERE (.+)").WithArgs(
					channelSettingId).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec("INSERT INTO channel_categories").
					WithArgs(channelSettingId, args.channel.Setting.Categories[0]).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Update(testCase.args.channel)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestChannelPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChannelPostgres(db)

	type args struct {
		apiId int
	}

	type mockBehavior func(args args)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				apiId: 1111,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec("DELETE FROM channels WHERE (.+)").
					WithArgs(args.apiId).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Not Found",
			args: args{
				apiId: 1111,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec("DELETE FROM channels WHERE (.+)").
					WithArgs(args.apiId).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err := r.Delete(testCase.args.apiId)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
