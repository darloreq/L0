package service

import (
	"encoding/json"
	"errors"
	"github.com/darloreq/L0/model"
)

type database interface {
	GetOrders() ([]model.Order, error) //Получение заказов из БД
	AddOrder(order model.Order) error  //Добавление заказа в БД
}

type cache interface {
	AddValue(order model.Order)
	GetValue(key string) (value string, bool2 bool)
}

type OrderService struct { //структура БЛ для обработки заказов
	Orderdb    database
	OrderCache cache
}

func (n OrderService) AddOrder(orderJson string) error { //Метод выполняет addorder с экземляром структуры n, получая строку от уровня представления, и возвращает ему либо ничего, либо ошибку
	order := model.Order{}
	err := json.Unmarshal([]byte(orderJson), &order)
	if err != nil {
		return err
	}
	order.OrderInfo = orderJson
	err = n.Orderdb.AddOrder(order) //запись в БД
	if err != nil {
		return err
	}
	n.OrderCache.AddValue(order) //запись в кеш
	return nil
}

func (n OrderService) GetOrder(uid string) (string, error) {
	orderInfo, ok := n.OrderCache.GetValue(uid)
	if !ok {
		return "", errors.New("произошла ошибка при получении данных из кеша, этих данных нет в кеше")
	}
	return orderInfo, nil
}

//БЛ от кеша хочет записывать данные и получать
//внедрить кеш в БЛ, доработать addorder, чтобы он записывал и в бд и в кеш, написать getorder, чтобы он выдавал данные из кеша по uid
