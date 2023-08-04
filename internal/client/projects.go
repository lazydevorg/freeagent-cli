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

func (c *Client) GetActiveProjects() ([]Project, error) {
	return GetCollection[Project](c, "projects?view=active", "projects", nil)
}
