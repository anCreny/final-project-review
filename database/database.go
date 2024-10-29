package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sater-151/todo-list/models"
	_ "modernc.org/sqlite"
)

type DBStruct struct {
	db *sql.DB
}

func CreateDB(dbFile string) error {
	db, _ := sql.Open("sqlite", dbFile)
	_, err := db.Exec(`CREATE TABLE scheduler (
    "id" INTEGER NOT NULL PRIMARY KEY,
    "date" INTEGER NOT NULL,
    "title" VARCHAR(100) NOT NULL DEFAULT "",
    "comment" VARCHAR(250),
    "repeat" VARCHAR(120) NOT NULL
	);`)
	if err != nil {
		log.Print("Ошибка создания БД: ", err.Error())
		return err
	}
	_, err = db.Exec("create index scheduler_data on scheduler (date);")
	if err != nil {
		log.Print("Ошибка создания индекса в БД: ", err.Error())
		return err
	}
	return nil
}

func OpenDB(DbFilePath string) (DBStruct, error) {
	var DB DBStruct

	dbFile := filepath.Join(DbFilePath, "scheduler.db")
	fmt.Println(dbFile)
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		err = CreateDB(dbFile)
		if err != nil {
			log.Print(err.Error())
			return DB, err
		}
	}
	DB.db, _ = sql.Open("sqlite", dbFile)
	return DB, nil
}

func (DB DBStruct) Close() {
	DB.db.Close()
}

func (DB DBStruct) InsertTask(task models.Task) (string, error) {
	res, err := DB.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}

func (DB DBStruct) SelectSortDate() ([]models.Task, error) {
	tasks := []models.Task{}
	res, err := DB.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT 20")
	if err != nil {
		return tasks, err
	}
	for res.Next() {
		task := models.Task{}
		check := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if check != nil {
			break
		}
		tasks = append(tasks, task)
	}
	defer res.Close()
	return tasks, nil
}

func (DB DBStruct) UpdateTask(task models.Task) error {
	res, err := DB.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("data update error")
	}
	return nil
}

func (DB DBStruct) DeleteTask(id string) error {
	_, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = DB.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func (DB DBStruct) Select(selectConfig models.SelectConfig) ([]models.Task, error) {
	listTask := []models.Task{}
	row := fmt.Sprintf("SELECT * FROM %s", selectConfig.Table)
	if selectConfig.Search != "" || selectConfig.Date != "" || selectConfig.Id != "" {
		row += " WHERE"
	}
	if selectConfig.Search != "" {
		row += fmt.Sprintf(" title LIKE %s OR comment LIKE %s", "'%"+selectConfig.Search+"%'", "'%"+selectConfig.Search+"%'")
	}
	if selectConfig.Date != "" {
		row += fmt.Sprintf(" date = '%s'", selectConfig.Date)
	}
	if selectConfig.Id != "" {
		row += fmt.Sprintf(" id = %s", selectConfig.Id)
	}
	if selectConfig.Sort != "" {
		row += fmt.Sprintf(" ORDER BY %s %s", selectConfig.Sort, selectConfig.TypeSort)
	}
	if selectConfig.Limit != "" {
		row += fmt.Sprintf(" LIMIT %s", selectConfig.Limit)
	}
	res, err := DB.db.Query(row)
	if err != nil {
		return listTask, err
	}
	for res.Next() {
		task := models.Task{}
		check := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if check != nil {
			break
		}
		listTask = append(listTask, task)
	}
	defer res.Close()
	return listTask, nil
}
