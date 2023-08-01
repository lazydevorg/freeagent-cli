package client

import "time"

type Project struct {
	Url       string    `json:"url"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (p Project) String() string {
	return p.Name
}

func GetActiveProjects() ([]Project, error) {
	return GetCollection[Project]("projects?view=active", "projects", nil)
}
