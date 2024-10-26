package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sater-151/todo-list/database"
	m "github.com/sater-151/todo-list/models"
	"github.com/sater-151/todo-list/services"
)

func GetNextDate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	now := req.FormValue("now")
	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	nextDate, err := services.NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(nextDate))
}

func PostTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var task m.Task
	var buf bytes.Buffer
	var errJS m.Error
	var idJS m.ID
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if task.Title == "" {
		errJS.Err = "Title is empty"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			errJS.Err = "Date is uncorrect form"
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(errJS)
			return
		}
		now := time.Now().Format("20060102")
		if t.Compare(time.Now()) == -1 && now != task.Date {
			if task.Repeat == "" {
				task.Date = time.Now().Format("20060102")
			} else {
				task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					errJS.Err = err.Error()
					res.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(res).Encode(errJS)
					return
				}
			}
		}
	}
	idJS.ID, err = database.InsertTask(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(idJS)
}

func GetTasks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var errJS m.Error
	search := req.FormValue("search")
	if search == "" {
		tasks, err := database.SelectSortDate()
		if err != nil {
			errJS.Err = err.Error()
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(errJS)
			return
		}
		tasksJS := m.TasksJS{Tasks: tasks}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(tasksJS)
		return
	} else {
		tasks, err := database.SelectBySearch(search)
		if err != nil {
			errJS.Err = err.Error()
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(errJS)
			return
		}
		tasksJS := m.TasksJS{Tasks: tasks}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(tasksJS)
		return
	}
}

func GetTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var errJS m.Error
	id := req.FormValue("id")
	if id == "" {
		errJS.Err = "uncorrect id"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	task, err := database.SelectByID(id)
	if err != nil {
		errJS.Err = "no task"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(task)
}

func PutTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var task m.Task
	var buf bytes.Buffer
	var errJS m.Error
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if task.Title == "" {
		errJS.Err = "Title is empty"
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			errJS.Err = "Date is uncorrect form"

			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(errJS)
			return
		}
		now := time.Now().Format("20060102")
		if t.Compare(time.Now()) == -1 && now != task.Date {
			if task.Repeat == "" {
				task.Date = time.Now().Format("20060102")
			} else {
				task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					errJS.Err = err.Error()
					res.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(res).Encode(errJS)
					return
				}
			}
		}
	}
	err = database.UpdateTask(task)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})

}

func PostTaskDone(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var errJS m.Error
	id := req.FormValue("id")
	task, err := database.SelectByID(id)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			errJS.Err = err.Error()
			res.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(res).Encode(errJS)
			return
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(struct{}{})
		return
	}

	task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	err = database.UpdateTask(task)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})
}

func DeleteTask(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	var errJS m.Error
	id := req.FormValue("id")
	err := database.DeleteTask(id)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(struct{}{})

}

func Sign(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json; charset=UTF-8")
	passTrue := os.Getenv("TODO_PASSWORD")
	var passJS m.PasswordJS
	var errJS m.Error
	var token m.JWTToken
	err := json.NewDecoder(req.Body).Decode(&passJS)
	if err != nil {
		errJS.Err = err.Error()
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	if passTrue != passJS.Pass {
		errJS.Err = "Неверный пароль"
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token.Token, err = jwtToken.SignedString([]byte(passTrue))
	if err != nil {
		errJS.Err = fmt.Sprintf("filed to sign jwt: %v", err.Error())
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(errJS)
		return
	}
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(token)
}

func Auth(n http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := req.Cookie("token")
			if err != nil {
				http.Error(res, err.Error(), http.StatusUnauthorized)
				return
			}
			jwtCookie := cookie.Value
			jwtToken, err := jwt.Parse(jwtCookie, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil || !jwtToken.Valid {
				http.Error(res, err.Error(), http.StatusUnauthorized)
				return
			}
		}
		n(res, req)
	})
}
