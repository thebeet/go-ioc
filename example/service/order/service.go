package order

import (
	"database/sql"
	"example/service/user"

	"github.com/thebeet/go-ioc/pkg/ioc"
)

func init() {
	ioc.Register(New)
}

type Service interface {
	GetOrderById(id int) *Order
}

func New() Service {
	return &service{}
}

type service struct {
	Db          *sql.DB      `autowire:"db_order"`
	UserService user.Service `autowire:""`
}

func (s *service) GetOrderById(id int) *Order {
	order := &Order{}
	var userName string
	oSql := "SELECT `id`, `user_name`, `item` FROM `order` WHERE `id`=? LIMIT 1"
	s.Db.QueryRow(oSql, id).Scan(&order.Id, &userName, &order.Item)
	order.User = *s.UserService.GetUserByName(userName)
	return order
}
