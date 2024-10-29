package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sater-151/todo-list/database"
	"github.com/sater-151/todo-list/handlers"
	"github.com/sater-151/todo-list/utils"
)

func main() {
	port, dbFilePath := utils.Config()

	Db, err := database.OpenDB(dbFilePath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer Db.Close()

	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Get("/api/tasks", handlers.Auth(handlers.ListTask(Db)))
	r.Get("/api/task", handlers.Auth(handlers.GetTask(Db)))

	r.Post("/api/task", handlers.Auth(handlers.PostTask(Db)))
	r.Post("/api/task/done", handlers.Auth(handlers.PostTaskDone(Db)))
	r.Post("/api/signin", handlers.Sign)

	r.Put("/api/task", handlers.Auth(handlers.PutTask(Db)))

	r.Delete("/api/task", handlers.Auth(handlers.DeleteTask(Db)))

	log.Println("Server start at port:", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Println("Ошибка запуска сервера:", err.Error())
		return
	}
}
