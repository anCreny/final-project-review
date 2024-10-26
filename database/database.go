package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	m "github.com/sater-151/todo-list/models"
	_ "modernc.org/sqlite"
)

type DBStruct struct {
	db *sql.DB
}

var DB DBStruct

func CreateDB() {
	_, err := DB.db.Exec(`CREATE TABLE scheduler (
    "id" INTEGER NOT NULL PRIMARY KEY,
    "date" integer NOT NULL,
    "title" VARCHAR(100) NOT NULL DEFAULT "",
    "comment" VARCHAR(250),
    "repeat" VARCHAR(120) NOT NULL
	);`)
	if err != nil {
		log.Print("Ошибка создания БД: ", err.Error())
		return
	}
	_, err = DB.db.Exec("create index scheduler_data on scheduler (date);")
	if err != nil {
		log.Print("Ошибка создания индекса в БД: ", err.Error())
		return
	}
}

func DBClose() {
	DB.db.Close()
}

func OpenDB(DbFilePath string) {
	dbFile := filepath.Join(DbFilePath, "scheduler.db")
	_, err := os.Stat(dbFile)
	db, _ := sql.Open("sqlite", dbFile)
	DB.db = db
	if os.IsNotExist(err) {
		CreateDB()
	}
}

func InsertTask(task m.Task) (string, error) {
	res, err := DB.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) values (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	id, _ := res.LastInsertId()
	return strconv.Itoa(int(id)), nil
}

func SelectSortDate() ([]m.Task, error) {
	tasks := []m.Task{}
	res, err := DB.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT 20")
	if err != nil {
		return tasks, err
	}
	for res.Next() {
		task := m.Task{}
		check := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if check != nil {
			break
		}
		tasks = append(tasks, task)
	}
	defer res.Close()
	return tasks, nil
}

func SelectBySearch(search string) ([]m.Task, error) {
	tasks := []m.Task{}
	var res *sql.Rows
	var d string
	var err error
	date := strings.Split(search, ".")
	if len(date) == 3 {
		for i := 2; i >= 0; i-- {

			d += date[i]
		}
		_, err := time.Parse("20060102", d)
		if err == nil {
			date, _ := strconv.Atoi(d)
			res, err = DB.db.Query("SELECT * FROM scheduler WHERE date = :d ORDER BY date LIMIT 20", sql.Named("d", date))
			if err != nil {
				return tasks, err
			}
		}
	} else {
		res, err = DB.db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT 20", sql.Named("search", "%"+search+"%"))
		if err != nil {
			return tasks, err
		}
	}
	for res.Next() {
		task := m.Task{}
		check := res.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if check != nil {
			break
		}
		tasks = append(tasks, task)
	}
	defer res.Close()
	return tasks, nil
}

func SelectByID(id string) (m.Task, error) {
	task := m.Task{}
	var date int
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return task, fmt.Errorf("uncorrect id")
	}
	res := DB.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	check := res.Scan(&idNum, &date, &task.Title, &task.Comment, &task.Repeat)
	if check != nil {
		return task, fmt.Errorf("Uncorrect id")
	}
	task.ID = strconv.Itoa(idNum)
	task.Date = strconv.Itoa(date)
	return task, nil
}

func UpdateTask(task m.Task) error {
	_, err := SelectByID(task.ID)
	if err != nil {
		return err
	}
	_, err = DB.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	return nil
}

func DeleteTask(id string) error {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = DB.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", idNum))
	if err != nil {
		return err
	}
	return nil
}
