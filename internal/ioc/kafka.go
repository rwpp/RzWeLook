package ioc

import (
	"github.com/IBM/sarama"
	"github.com/rwpp/RzWeLook/internal/events"
	"github.com/rwpp/RzWeLook/internal/events/article"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) sarama.SyncProducer {
	//saramaCfg := sarama.NewConfig()
	//saramaCfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}

func NewConsumers(client *article.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{client}
}
