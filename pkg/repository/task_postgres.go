package repository

import (
	"fmt"
	"strings"
	"time"

	core "github.com/max-sanch/BotFreelancer-core"

	"github.com/jmoiron/sqlx"
)

type TaskPostgres struct {
	db *sqlx.DB
}

func NewTaskPostgres(db *sqlx.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) GetOrCreateCategoryByName(name string) (int, error) {
	var id int

	getQuery := fmt.Sprintf("SELECT id FROM %s WHERE lower(name) = $1;", categoriesTable)

	row := r.db.QueryRow(getQuery, strings.ToLower(name))
	if err := row.Scan(&id); err != nil {
		createQuery := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id;", categoriesTable)
		row = r.db.QueryRow(createQuery, name)
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (r *TaskPostgres) GetLastParseTime() (string, error) {
	var datetime string

	query := fmt.Sprintf("SELECT datetime FROM %s WHERE id = 1;", lastParsedTasksTable)
	row := r.db.QueryRow(query)
	if err := row.Scan(&datetime); err != nil {
		return "", nil
	}

	return datetime, nil
}

func (r *TaskPostgres) SetLastParseTime() error {
	currentTime := time.Now().UTC()

	query := fmt.Sprintf(`
		INSERT INTO %s (id, datetime) VALUES
		(1, TIMESTAMP WITH TIME ZONE '%s+00') ON CONFLICT (id)
		DO UPDATE SET datetime = TIMESTAMP WITH TIME ZONE '%s+00'
		WHERE %s.id = 1;`,
		lastParsedTasksTable, currentTime.Format("2006-01-02 15:04:05"),
		currentTime.Format("2006-01-02 15:04:05"), lastParsedTasksTable)

	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) GetAllForChannels() ([]core.ChannelTaskResponse, error) {
	var tasks []core.ChannelTaskResponse

	query := fmt.Sprintf(`SELECT ch.api_id, ch.api_hash, flt.title, flt.body, flt.task_url FROM %s ch
		INNER JOIN %s chs ON ch.id = chs.channel_id
		INNER JOIN %s flt ON flt.is_budget = chs.is_budget AND flt.is_term = chs.is_term AND
		flt.is_safe_deal = chs.is_safe_deal AND
		flt.category_id in (SELECT category_id FROM %s WHERE channel_setting_id = chs.id)
		ORDER BY ch.id, flt.id;`,
		channelsTable, channelSettingsTable, freelanceTasksTable, channelCategoriesTable)

	if err := r.db.Select(&tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskPostgres) GetAllForUsers() ([]core.UserTaskResponse, error) {
	var tasks []core.UserTaskResponse

	query := fmt.Sprintf(`SELECT u.tg_id, flt.title, flt.body, flt.task_url FROM %s u
		INNER JOIN %s us ON u.id = us.user_id
		INNER JOIN %s flt ON flt.is_budget = us.is_budget AND flt.is_term = us.is_term AND
		flt.is_safe_deal = us.is_safe_deal AND
		flt.category_id in (SELECT category_id FROM %s WHERE user_setting_id = us.id)
		ORDER BY u.id, flt.id;`,
		usersTable, userSettingsTable, freelanceTasksTable, userCategoriesTable)

	if err := r.db.Select(&tasks, query); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskPostgres) AddTasks(tasksInput core.TasksInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, task := range tasksInput.Tasks {
		isBudget := task.Budget != 0
		isTerm := task.Term != ""

		categoryId, err := r.GetOrCreateCategoryByName(task.Category)
		if err != nil {
			return err
		}

		body, err := getTaskBody(task, isBudget, isTerm)
		if err != nil {
			return err
		}

		createTaskQuery := fmt.Sprintf(`INSERT INTO %s (task_url, title, body, category_id, is_budget, is_term, is_safe_deal)
			VALUES ($1, $2, $3, $4, $5, $6, $7);`, freelanceTasksTable)

		if _, err := tx.Exec(createTaskQuery, task.TaskUrl, task.Title, body,
			categoryId, isBudget, isTerm, task.IsSafeDeal); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *TaskPostgres) DeleteAll() error {
	query := fmt.Sprintf("DELETE FROM %s;", freelanceTasksTable)
	if _, err := r.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func getTaskBody(task core.TaskDataInput, isBudget, isTerm bool) (string, error) {
	var body string
	var budget string
	term := "не указаны"
	safeDeal := ""

	if isBudget {
		budget = fmt.Sprintf("%d ", task.Budget)
		if task.IsBudgetPerHour {
			budget += "в час"
		}
	} else {
		budget = "не указан"
	}

	if isTerm {
		term = task.Term
	}

	if task.IsSafeDeal {
		safeDeal = "Безопасная сделка!\n"
	}

	body = fmt.Sprintf("Заказ с %s\n\nКатегория: %s\n\nОписание:\n%s\n\n%sБюджет: %s\nСроки: %s\nВремя публикации: %s",
		task.FLName, task.Category, task.Description, safeDeal, budget, term, task.DateTime)

	return body, nil
}
