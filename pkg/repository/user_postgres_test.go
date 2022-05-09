package repository

import (
	"database/sql"
	"errors"
	"testing"

	core "github.com/max-sanch/BotFreelancer-core"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserPostgres_GetByTgId(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPostgres(db)

	type args struct {
		tgId int
	}

	type mockBehavior func(args args)

	testTables := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		want         core.UserResponse
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				tgId: 1111,
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "tg_id", "username"}).
					AddRow(1, 1111, "user-1")

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(args.tgId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id", "is_safe_deal", "is_budget", "is_term"}).
					AddRow(1, true, true, true)

				mock.ExpectQuery("SELECT (.+) FROM user_settings WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"category_id"}).AddRow(1).AddRow(2)

				mock.ExpectQuery("SELECT (.+) FROM user_categories WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			want: core.UserResponse{
				Id:       1,
				TgId:     1111,
				Username: "user-1",
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
				tgId: 1111,
			},
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "tg_id", "username"})

				mock.ExpectQuery("SELECT (.+) FROM users WHERE (.+)").
					WithArgs(args.tgId).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTables {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			got, err := r.GetByTgId(testCase.args.tgId)
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

func TestUserPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPostgres(db)
	isFalse := false

	type args struct {
		user core.UserInput
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
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					args.user.TgId, args.user.Username).WillReturnRows(rows)

				userSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(userSettingId)
				mock.ExpectQuery("INSERT INTO user_settings").WithArgs(
					id, args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm).WillReturnRows(rows)

				for _, categoryId := range args.user.Setting.Categories {
					mock.ExpectExec("INSERT INTO user_categories").WithArgs(
						userSettingId, categoryId).WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "",
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
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					args.user.TgId, args.user.Username).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					args.user.TgId, args.user.Username).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO user_settings").WithArgs(
					id, args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "3nd Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("INSERT INTO users").WithArgs(
					args.user.TgId, args.user.Username).WillReturnRows(rows)

				userSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(userSettingId)
				mock.ExpectQuery("INSERT INTO user_settings").WithArgs(
					id, args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO user_categories").
					WithArgs(userSettingId, args.user.Setting.Categories[0]).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}

func TestUserPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserPostgres(db)
	isFalse := false

	type args struct {
		user core.UserInput
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
				user: core.UserInput{
					TgId:     1111,
					Username: "channel-1",
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
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WithArgs(
					args.user.Username, args.user.TgId).WillReturnRows(rows)

				userSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(userSettingId)
				mock.ExpectQuery("UPDATE user_settings SET (.+) WHERE (.+)").WithArgs(
					args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM user_categories WHERE (.+)").WithArgs(
					userSettingId).WillReturnResult(sqlmock.NewResult(0, 1))

				for _, categoryId := range args.user.Setting.Categories {
					mock.ExpectExec("INSERT INTO user_categories").WithArgs(
						userSettingId, categoryId).WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "",
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
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WithArgs(
					args.user.Username, args.user.TgId).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WithArgs(
					args.user.Username, args.user.TgId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("UPDATE user_settings SET (.+) WHERE (.+)").WithArgs(
					args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Not Found",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WithArgs(
					args.user.Username, args.user.TgId).WillReturnRows(rows)

				userSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(userSettingId)
				mock.ExpectQuery("UPDATE user_settings SET (.+) WHERE (.+)").WithArgs(
					args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM user_categories WHERE (.+)").WithArgs(
					userSettingId).WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "3nd Empty Fields",
			args: args{
				user: core.UserInput{
					TgId:     1111,
					Username: "user-1",
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
				mock.ExpectQuery("UPDATE users SET (.+) WHERE (.+)").WithArgs(
					args.user.Username, args.user.TgId).WillReturnRows(rows)

				userSettingId := 3
				rows = sqlmock.NewRows([]string{"id"}).AddRow(userSettingId)
				mock.ExpectQuery("UPDATE user_settings SET (.+) WHERE (.+)").WithArgs(
					args.user.Setting.IsSafeDeal, args.user.Setting.IsBudget,
					args.user.Setting.IsTerm, id).WillReturnRows(rows)

				mock.ExpectExec("DELETE FROM user_categories WHERE (.+)").WithArgs(
					userSettingId).WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectExec("INSERT INTO user_categories").
					WithArgs(userSettingId, args.user.Setting.Categories[0]).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Update(testCase.args.user)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}
}
