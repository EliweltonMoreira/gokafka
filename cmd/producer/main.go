package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	err := Publish("transferiu", "teste", producer, []byte("transferencia2"), deliveryChan)
	if err != nil {
		log.Println(err.Error())
	}
	DeliveryReport(deliveryChan) // async

	//e := <-deliveryChan
	//msg := e.(*kafka.Message)
	//if msg.TopicPartition.Error != nil {
	//	fmt.Println("Erro ao enviar")
	//} else {
	//	fmt.Println("Mensagem enviada:", msg.TopicPartition)
	//}
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "gokafka-kafka-1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}
	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada:", ev.TopicPartition)
				// anotar no banco de dados que a mensagem foi processada.
				// ex: confirma que uma transferência bancária ocorreu.
			}
		}
	}
}
