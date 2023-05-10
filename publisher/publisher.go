package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"time"
)

func main() {
	fu := func(msg *stan.Msg) {
		fmt.Println(string(msg.Data))
	}
	models := []string{"../model.json", "../model1.json", "../model2.json", "../model4.json"}
	sc, err := stan.Connect("test-cluster", "1234")
	if err != nil {
		fmt.Println(err, sc)
		return
	}
	defer sc.Close()
	subscribe, err := sc.Subscribe("1231", fu, stan.StartAtTimeDelta(time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer subscribe.Close()
	for _, str := range models {
		file, err := ioutil.ReadFile(str)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = sc.Publish("1231", file)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Message sent!")
		}
	}
}
