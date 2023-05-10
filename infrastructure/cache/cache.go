package cache

import (
	"github.com/darloreq/L0/model"
	"sync"
)

type cache struct {
	cachemap map[string]string
	mutex    *sync.Mutex
}

type database interface { //интерфейс, необходимый для обращения к БД для кеша
	GetOrders() ([]model.Order, error)
}

func New(d database) (*cache, error) { //функция, позволяющая добавить в кеш данные из БД, при рестарте сервиса
	retCache := &cache{
		cachemap: make(map[string]string),
		mutex:    &sync.Mutex{},
	}
	orders, err := d.GetOrders()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(orders); i++ {
		retCache.AddValue(orders[i])
	}
	return retCache, nil
}

func (c cache) AddValue(order model.Order) { //метод добавления в кеш
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cachemap[order.Uid] = order.OrderInfo
}

func (c cache) GetValue(key string) (value string, bool2 bool) { //метод взятия данных из кеша
	c.mutex.Lock()
	defer c.mutex.Unlock()
	value, ok := c.cachemap[key]
	if ok != true {
		return value, false
	}
	return value, true
}
