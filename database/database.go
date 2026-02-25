package database

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Ctx  context.Context
	Conn *pgx.Conn
}

type TaskModel struct {
	Id         int
	Title      string
	Desc       string
	Priority   int
	Status     bool
	DateCreate time.Time
	DateCompl  *time.Time
}

func Connect(ctx context.Context) (*pgx.Conn, error) {

	conn_string := os.Getenv("CONN_STRING")
	return pgx.Connect(ctx, conn_string)

}

func (d *Database) InsertConn(ctx context.Context, conn *pgx.Conn) {

	d.Ctx = ctx
	d.Conn = conn

}

func CreateTable(ctx context.Context, conn *pgx.Conn) error {

	sqlQuery := `
	CREATE TABLE IF NOT EXISTS tasks(
		id SERIAL PRIMARY KEY, 
		title VARCHAR(50) NOT NULL, 
		description VARCHAR(1000), 
		priority INT NOT NULL, 
		status BOOLEAN NOT NULL, 
		date_create TIMESTAMP NOT NULL, 
		date_completed TIMESTAMP,
		UNIQUE(title)
	);
	`
	_, err := conn.Exec(ctx, sqlQuery)
	return err

}

func (db *Database) Insert(ctx context.Context, conn *pgx.Conn, title, description string, priority int, status bool, date_create time.Time) error {

	sqlQuery := `INSERT INTO tasks (title, description, priority, status, date_create)
	VALUES ($1, $2, $3, $4, $5);`
	_, err := conn.Exec(ctx, sqlQuery, title, description, priority, status, date_create)
	return err

}

func (db *Database) Delete(ctx context.Context, conn *pgx.Conn, title string) error {

	sqlQuery := `DELETE FROM tasks WHERE title = $1;`
	_, err := conn.Exec(ctx, sqlQuery, title)
	return err

}

func (db *Database) DeleteCompletedTasks(ctx context.Context, conn *pgx.Conn) error {

	sqlQuery := `DELETE FROM tasks WHERE status = TRUE;`
	_, err := conn.Exec(ctx, sqlQuery)
	return err

}

func (db *Database) Update(ctx context.Context, conn *pgx.Conn, task TaskModel, title string) error {

	var complTime interface{}
	if task.DateCompl != nil {
		complTime = *task.DateCompl
	} else {
		complTime = nil
	}

	sqlQuery := `UPDATE tasks SET title = $1, description = $2, priority = $3, status = $4, date_completed = $5 WHERE title = $6;`
	_, err := conn.Exec(ctx, sqlQuery, task.Title, task.Desc, task.Priority, task.Status, complTime, title)
	return err

}

func (db Database) GetTask(ctx context.Context, conn *pgx.Conn, title string) (error, TaskModel) {

	task := TaskModel{}
	sqlQuery := `SELECT * FROM tasks WHERE title = $1;`
	err := conn.QueryRow(ctx, sqlQuery, title).Scan(&task.Id, &task.Title, &task.Desc, &task.Priority, &task.Status, &task.DateCreate, &task.DateCompl)
	return err, task

}

func (db Database) GetAllTasks(ctx context.Context, conn *pgx.Conn) (error, []TaskModel) {

	tasks := []TaskModel{}
	task := TaskModel{}
	sqlQuery := `SELECT * FROM tasks;`
	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&task.Id, &task.Title, &task.Desc, &task.Priority, &task.Status, &task.DateCreate, &task.DateCompl); err != nil {
			panic(err)
		}
		tasks = append(tasks, task)
	}
	return err, tasks

}
