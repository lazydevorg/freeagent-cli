package client

import (
	"fmt"
	"time"
)

type User struct {
	Url             string    `json:"url"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Role            string    `json:"role"`
	PermissionLevel int       `json:"permission_level"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (u User) String() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (c *Client) PersonalProfile() (*User, error) {
	return GetEntity[User](c, "users/me", "user")
}
