package kafkaapi

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/coreos/pkg/capnslog"
	"reflect"
	"staticzeng.com/config"
)

var (
	kp    *kafka.Producer
	topic string
	log   = capnslog.NewPackageLogger(reflect.TypeOf(struct{}{}).PkgPath(), "kafkaapi.kafkaapi")
)

func init() {
	conf := config.LoadConfig()
	log.Info("conf.Broker : ", conf.Kafka.Broker)
	log.Info("conf.Topic : ", conf.Kafka.Topic)
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": conf.Kafka.Broker})
	if err != nil {
		log.Info("create producer error : ", err)
	} else {
		kp = producer
	}
	topic = conf.Kafka.Topic
}

func Produce(msg string) {
	kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg),
	}, nil)
}
