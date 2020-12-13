package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"iotflux/config"
	"iotflux/middleware/kafka"
)
var cfg = new (config.IotfluxConfig)

func main(){

    err:= ini.MapTo(cfg, "./config/config.ini")
    if err != nil{
    	fmt.Println("load ini failed err", err)
		return
	}

	//加载配置文件
	err = kafka.InitProducer(cfg.KafkaConfig.Address)
	if err != nil{
		fmt.Println("Initing kafka Error:", err)
		return
	}
	for i:=0; i<20; i++ {
		go kafka.SendMessage(cfg.Topic, "I Love You")

	}

	consumer, err := kafka.NewKafkaConsumer(cfg.KafkaConfig.Address, nil )
	kafka.ReadMessage(consumer, cfg.KafkaConfig.Topic)

}