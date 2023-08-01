package client

import "time"

type Task struct {
	Url       string    `json:"url"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (t Task) String() string {
	return t.Name
}

func GetActiveTasks() ([]Task, error) {
	return GetCollection[Task]("tasks?view=active", "tasks", nil)
}
