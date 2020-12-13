package config

type KafkaConfig struct {
	Address []string `ini:"address"`
	Topic string `ini:"topic"`
}


type IotfluxConfig struct {
	KafkaConfig `ini:"Kafka"`
}