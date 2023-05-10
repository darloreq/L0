package repository

import (
	"database/sql"
	"fmt"
	"github.com/darloreq/L0/model"
)

type database struct {
	db *sql.DB
}

func New(db *sql.DB) *database {
	return &database{db: db}
}

func (d database) GetOrders() ([]model.Order, error) { //метод, описывающий получение данных из бд, использующий экземпляр
	rows, err := d.db.Query("select * from orders")
	if err != nil {
		return nil, err
	}
	orders := make([]model.Order, 0)
	for rows.Next() {
		o := model.Order{}
		err := rows.Scan(&o.Uid, &o.OrderInfo)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (d database) AddOrder(order model.Order) error { //метод, описывающий добавление(запись) данных из структуры ордер в БД
	_, err := d.db.Exec("insert into orders(uid, order_info) values ($1, $2)", order.Uid, order.OrderInfo)
	if err != nil {
		return err
	}
	return nil
}
