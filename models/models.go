package models

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ID struct {
	ID string `json:"id"`
}

type Error struct {
	Err string `json:"error"`
}

type TasksJS struct {
	Tasks []Task `json:"tasks"`
}

type PasswordJS struct {
	Pass string `json:"password"`
}

type JWTToken struct {
	Token string `json:"token"`
}
