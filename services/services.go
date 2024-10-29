package services

import (
	"strings"
	"time"

	"github.com/sater-151/todo-list/database"
	"github.com/sater-151/todo-list/models"
	"github.com/sater-151/todo-list/utils"
)

func AddTask(db database.DBStruct, task models.Task) (models.ID, error) {
	var Id models.ID
	var err error
	task, err = utils.CheckTask(task)
	if err != nil {
		return Id, err
	}
	Id.ID, err = db.InsertTask(task)
	if err != nil {
		return Id, err
	}
	return Id, nil
}

func UpdateTask(db database.DBStruct, task models.Task) error {
	task, err := utils.CheckTask(task)
	if err != nil {
		return err
	}
	err = db.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func TaskDone(db database.DBStruct, selectConfig models.SelectConfig) error {
	tasks, err := db.Select(selectConfig)
	if err != nil {
		return err
	}
	task := tasks[0]
	if task.Repeat == "" {
		err = db.DeleteTask(selectConfig.Id)
		if err != nil {
			return err
		}
		return nil
	}
	task.Date, err = utils.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}
	err = db.UpdateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func GetListTask(db database.DBStruct, selectConfig models.SelectConfig) ([]models.Task, error) {
	if selectConfig.Search != "" {
		date := strings.Split(selectConfig.Search, ".")
		if len(date) == 3 {
			var d string
			for i := 2; i >= 0; i-- {
				d += date[i]
			}
			_, err := time.Parse("20060102", d)
			if err == nil {
				selectConfig.Search = ""
				selectConfig.Date = d
			}
		}
	}
	listTask, err := db.Select(selectConfig)
	if err != nil {
		return listTask, err
	}
	return listTask, err
}
