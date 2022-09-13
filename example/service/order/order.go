package order

import "example/service/user"

type Order struct {
	Id   int       `json:"id"`
	User user.User `json:"user"`
	Item string    `json:"item"`
}
