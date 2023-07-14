package client

import (
	"time"
)

type UsersClient struct {
	client *Client
}

func (c *Client) NewUsersClient() *UsersClient {
	return &UsersClient{client: c}
}

type User struct {
	Url                string    `json:"url"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	PermissionLevel    int       `json:"permission_level"`
	NiNumber           string    `json:"ni_number"`
	UniqueTaxReference string    `json:"unique_tax_reference"`
	OpeningMileage     string    `json:"opening_mileage"`
	UpdatedAt          time.Time `json:"updated_at"`
	CreatedAt          time.Time `json:"created_at"`
}

func (u *UsersClient) PersonalProfile() (user *User, err error) {
	user, err = GetEntity[User](u.client, "https://api.sandbox.freeagent.com/v2/users/me", "user")
	return
}

func (u *UsersClient) GetAllUsers() (users []User, err error) {
	users, err = GetArray[User](u.client, "https://api.sandbox.freeagent.com/v2/users", "users")
	return
}
