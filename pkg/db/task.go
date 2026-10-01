package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	result, err := database.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat)
		 VALUES (?, ?, ?, ?)`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	task := &Task{}

	err := database.QueryRow(
		`SELECT id, date, title, comment, repeat
		 FROM scheduler
		 WHERE id = ?`,
		id,
	).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func DeleteTask(id string) error {
	result, err := database.Exec(
		`DELETE FROM scheduler WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateTask(task *Task) error {
	result, err := database.Exec(
		`UPDATE scheduler
         SET date = ?, title = ?, comment = ?, repeat = ?
         WHERE id = ?`,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
func Tasks(limit int, search string, date string) ([]*Task, error) {
	var (
		query string
		args  []any
	)

	switch {
	case date != "":
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE date = ?
			ORDER BY date
			LIMIT ?
		`
		args = []any{date, limit}

	case search != "":
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			WHERE title LIKE ? OR comment LIKE ?
			ORDER BY date
			LIMIT ?
		`

		searchPattern := "%" + search + "%"
		args = []any{searchPattern, searchPattern, limit}

	default:
		query = `
			SELECT id, date, title, comment, repeat
			FROM scheduler
			ORDER BY date
			LIMIT ?
		`
		args = []any{limit}
	}

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Проверяем, что sql импорт действительно используется.
// Это можно удалить, если он больше нигде не нужен.
var _ = sql.ErrNoRows
