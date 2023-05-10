package bootstrap

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) { //Метод, подключающий нашу БД для работы в го
	connStr := "user=postgres password=123 dbname=orders sslmode=disable" //Данные для подключения отправляются в sql
	db, err := sql.Open("postgres", connStr)                              //метод опен для инициализации экземпляра нашей бд и возможной ошибки при подключении
	if err != nil {
		return db, err
	}
	err = db.Ping() //чтобы понять, что подключение к бд есть и оно успешно, пингуем бд и проверяем на ошибку
	if err != nil {
		return nil, err
	}
	return db, err
}
