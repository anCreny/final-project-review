package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sater-151/todo-list/database"
	"github.com/sater-151/todo-list/handlers"
)

func setEnv() {
	port := "7540"
	dbFile := "/app"
	pass := "TestPas"
	err := os.Setenv("TODO_DBFILE", dbFile)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Setenv("TODO_PORT", port)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Setenv("TODO_PASSWORD", pass)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// setEnv()
	port := os.Getenv("TODO_PORT")
	dbFilePath := os.Getenv("TODO_DBFILE")

	database.OpenDB(dbFilePath)
	defer database.DBClose()

	r := chi.NewRouter()

	webDir := "web"
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	r.Get("/api/nextdate", handlers.GetNextDate)
	r.Get("/api/tasks", handlers.Auth(handlers.GetTasks))
	r.Get("/api/task", handlers.Auth(handlers.GetTask))

	r.Post("/api/task", handlers.Auth(handlers.PostTask))
	r.Post("/api/task/done", handlers.Auth(handlers.PostTaskDone))
	r.Post("/api/signin", handlers.Sign)

	r.Put("/api/task", handlers.Auth(handlers.PutTask))

	r.Delete("/api/task", handlers.Auth(handlers.DeleteTask))

	log.Printf("Server start at port: %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Print("Ошибка запуска сервера:", err.Error())
		return
	}
}
