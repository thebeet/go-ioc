package user

import (
	"database/sql"
	"log"

	"github.com/thebeet/go-ioc/pkg/ioc"
)

func init() {
	ioc.Register(New)
}

type Service interface {
	GetUserByName(name string) *User
}

func New(db *sql.DB) Service {
	return &service{db}
}

type service struct {
	Db *sql.DB
}

func (s *service) GetUserByName(name string) *User {
	user := &User{}
	uSql := "SELECT `id`, `name` FROM `user` WHERE `name` LIKE ? LIMIT 1"
	s.Db.QueryRow(uSql, name).Scan(&user.Id, &user.Name)
	log.Printf("get user %+v by name %s", user, name)
	return user
}
