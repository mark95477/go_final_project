package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func AddTask(t *Task) (string, error) {
	q := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?);"
	res, err := DB.Exec(q, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func Tasks(limit int) ([]*Task, error) {
	q := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?;"
	rows, err := DB.Query(q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*Task
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)
		list = append(list, &t)
	}
	return list, rows.Err()
}

func GetTask(id string) (*Task, error) {
	q := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?;"
	row := DB.QueryRow(q, id)

	var t Task
	var idInt int64
	err := row.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	t.ID = strconv.FormatInt(idInt, 10)
	return &t, nil
}

func UpdateTask(t *Task) error {
	q := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?;"
	res, err := DB.Exec(q, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found for update")
	}
	return nil
}

func DeleteTask(id string) error {
	q := "DELETE FROM scheduler WHERE id = ?;"
	res, err := DB.Exec(q, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found for deletion")
	}
	return nil
}

func UpdateDate(next, id string) error {
	q := "UPDATE scheduler SET date = ? WHERE id = ?;"
	res, err := DB.Exec(q, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found for date update")
	}
	return nil
}
