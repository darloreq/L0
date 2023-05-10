package main

import (
	"github.com/darloreq/L0/bootstrap"
	"github.com/darloreq/L0/infrastructure/cache"
	"github.com/darloreq/L0/infrastructure/repository"
	"github.com/darloreq/L0/model"
	"github.com/darloreq/L0/service"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"html/template"
	"log"
	"net/http"
)

func main() {
	db, _ := bootstrap.ConnectDB()
	q := repository.New(db)
	c, err := cache.New(q)
	if err != nil {
		log.Println(err)
		return
	}
	n := service.OrderService{Orderdb: q, OrderCache: c}

	// Подключение к nats-streaming
	connect, err := stan.Connect("test-cluster", "123")
	if err != nil {
		log.Println(err)
		return
	}
	defer connect.Close()
	// Анонимная функция обработки полученных данных из nats-streaming
	fu := func(msg *stan.Msg) {
		dataJson := string(msg.Data)
		err := n.AddOrder(dataJson)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Данные успешно добавлены!")
	}
	// Подписка/объявление канала nats-streaming
	subscribe, err := connect.Subscribe("1231", fu)
	if err != nil {
		log.Println(err)
		return
	}
	defer subscribe.Close()

	r := mux.NewRouter()
	// хендлеры сервера
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "htmlTemplates/FoundUid.html")
	})
	r.HandleFunc("/Uid", func(w http.ResponseWriter, r *http.Request) {

		uid := r.FormValue("uid")

		forUid, err := n.GetOrder(uid)
		order := model.Order{OrderInfo: forUid}
		tmpl, _ := template.ParseFiles("htmlTemplates/Uid.html")
		if err != nil {
			tmpl.Execute(w, order)
			return
		}
		tmpl.Execute(w, order)
	})

	http.Handle("/", r)
	http.ListenAndServe(":8181", nil)

}
