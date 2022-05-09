package repository

import (
	"errors"
	"fmt"
	"testing"

	core "github.com/max-sanch/BotFreelancer-core"

	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTaskPostgres_GetOrCreateCategoryByName(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPostgres(db)

	type args struct {
		name string
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
				name: "Category",
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO categories").WithArgs(args.name).WillReturnRows(rows)
			},
		},
		{
			name: "Empty Fields",
			args: args{
				name: "",
			},
			mockBehavior: func(args args, id int) {
				rows := sqlmock.NewRows([]string{"id"}).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO categories").WithArgs(args.name).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.GetOrCreateCategoryByName(testCase.args.name)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskPostgres_GetLastParseTime(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPostgres(db)

	testTable := []struct {
		name         string
		mockBehavior func()
		want         string
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"datetime"}).AddRow("datetime_test")
				mock.ExpectQuery("SELECT (.+) FROM last_parsed_tasks WHERE (.+)").WillReturnRows(rows)
			},
			want: "datetime_test",
		},
		{
			name: "Not Found",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"datetime"})
				mock.ExpectQuery("SELECT (.+) FROM last_parsed_tasks WHERE (.+)").WillReturnRows(rows)
			},
			want: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetLastParseTime()
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

func TestTaskPostgres_GetAllForChannels(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPostgres(db)

	testTable := []struct {
		name         string
		mockBehavior func()
		want         []core.ChannelTaskResponse
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"api_id", "api_hash", "title", "body", "task_url"}).
					AddRow(1111, "hash1111", "test", "test-body", "test-url").
					AddRow(1111, "hash1111", "test2", "test-body2", "test-url2").
					AddRow(3333, "hash3333", "test", "test-body", "test-url")
				mock.ExpectQuery("SELECT (.+) FROM channels ch INNER JOIN channel_settings chs ON (.+) INNER JOIN freelance_tasks flt ON (.+)").
					WillReturnRows(rows)
			},
			want: []core.ChannelTaskResponse{
				{
					ApiId:   1111,
					ApiHash: "hash1111",
					Title:   "test",
					Body:    "test-body",
					Url:     "test-url",
				},
				{
					ApiId:   1111,
					ApiHash: "hash1111",
					Title:   "test2",
					Body:    "test-body2",
					Url:     "test-url2",
				},
				{
					ApiId:   3333,
					ApiHash: "hash3333",
					Title:   "test",
					Body:    "test-body",
					Url:     "test-url",
				},
			},
		},
		{
			name: "Not Found",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"api_id", "api_hash", "title", "body", "task_url"})
				mock.ExpectQuery("SELECT (.+) FROM channels ch INNER JOIN channel_settings chs ON (.+) INNER JOIN freelance_tasks flt ON (.+)").
					WillReturnRows(rows)
			},
			want: []core.ChannelTaskResponse(nil),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetAllForChannels()
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

func TestTaskPostgres_GetAllForUsers(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPostgres(db)

	testTable := []struct {
		name         string
		mockBehavior func()
		want         []core.UserTaskResponse
		wantErr      bool
	}{
		{
			name: "OK",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"tg_id", "title", "body", "task_url"}).
					AddRow(1111, "test", "test-body", "test-url").
					AddRow(1111, "test2", "test-body2", "test-url2").
					AddRow(3333, "test", "test-body", "test-url")
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN user_settings us ON (.+) INNER JOIN freelance_tasks flt ON (.+)").
					WillReturnRows(rows)
			},
			want: []core.UserTaskResponse{
				{
					TgId:  1111,
					Title: "test",
					Body:  "test-body",
					Url:   "test-url",
				},
				{
					TgId:  1111,
					Title: "test2",
					Body:  "test-body2",
					Url:   "test-url2",
				},
				{
					TgId:  3333,
					Title: "test",
					Body:  "test-body",
					Url:   "test-url",
				},
			},
		},
		{
			name: "Not Found",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"tg_id", "title", "body", "task_url"})
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN user_settings us ON (.+) INNER JOIN freelance_tasks flt ON (.+)").
					WillReturnRows(rows)
			},
			want: []core.UserTaskResponse(nil),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			got, err := r.GetAllForUsers()
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

func TestTaskPostgres_AddTasks(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTaskPostgres(db)

	type args struct {
		tasksInput core.TasksInput
		categoryId int
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
				tasksInput: core.TasksInput{
					Tasks: []core.TaskDataInput{
						{
							FLName:          "test-FLName",
							FLUrl:           "test-FLUrl",
							TaskUrl:         "test-TaskUrl",
							Category:        "Category",
							Title:           "test-Title",
							Description:     "test-Description",
							Budget:          1000,
							IsBudgetPerHour: false,
							IsSafeDeal:      true,
							DateTime:        "test-DateTime",
						},
						{
							FLName:          "test-FLName2",
							FLUrl:           "test-FLUrl2",
							TaskUrl:         "test-TaskUrl2",
							Category:        "Category",
							Title:           "test-Title2",
							Description:     "test-Description2",
							Budget:          2000,
							IsBudgetPerHour: false,
							IsSafeDeal:      false,
							DateTime:        "test-DateTime2",
						},
					},
				},
				categoryId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				for _, task := range args.tasksInput.Tasks {
					rows := sqlmock.NewRows([]string{"id"}).AddRow(args.categoryId)
					mock.ExpectQuery("INSERT INTO categories").
						WillReturnRows(rows)

					budget := fmt.Sprintf("%d ", task.Budget)
					term := "не указаны"
					safeDeal := ""

					if task.IsSafeDeal {
						safeDeal = "Безопасная сделка!\n"
					}

					body := fmt.Sprintf("Заказ с %s\n\nКатегория: %s\n\nОписание:\n%s\n\n%sБюджет: %s\nСроки: %s\nВремя публикации: %s",
						task.FLName, task.Category, task.Description, safeDeal, budget, term, task.DateTime)

					mock.ExpectExec("INSERT INTO freelance_tasks").
						WithArgs(task.TaskUrl, task.Title, body, args.categoryId, true, false, task.IsSafeDeal).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				tasksInput: core.TasksInput{
					Tasks: []core.TaskDataInput{
						{
							FLName:          "",
							FLUrl:           "test-FLUrl",
							TaskUrl:         "test-TaskUrl",
							Category:        "Category",
							Title:           "test-Title",
							Description:     "test-Description",
							Budget:          1000,
							IsBudgetPerHour: false,
							IsSafeDeal:      true,
							DateTime:        "test-DateTime",
						},
					},
				},
				categoryId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectBegin()

				for _, task := range args.tasksInput.Tasks {
					rows := sqlmock.NewRows([]string{"id"}).AddRow(args.categoryId)
					mock.ExpectQuery("INSERT INTO categories").
						WillReturnRows(rows)

					budget := fmt.Sprintf("%d ", task.Budget)
					term := "не указаны"
					safeDeal := ""

					if task.IsSafeDeal {
						safeDeal = "Безопасная сделка!\n"
					}

					body := fmt.Sprintf("Заказ с %s\n\nКатегория: %s\n\nОписание:\n%s\n\n%sБюджет: %s\nСроки: %s\nВремя публикации: %s",
						task.FLName, task.Category, task.Description, safeDeal, budget, term, task.DateTime)

					mock.ExpectExec("INSERT INTO freelance_tasks").
						WithArgs(task.TaskUrl, task.Title, body, args.categoryId, true, false, task.IsSafeDeal).
						WillReturnError(errors.New("some error"))
				}

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			err := r.AddTasks(testCase.args.tasksInput)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
