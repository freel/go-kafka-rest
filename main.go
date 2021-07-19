package main

import (
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"go-kafka-rest/config"
	"log"
	"net/http"
)

var c config.RestConfig

func SendToKafka(w http.ResponseWriter, req *http.Request) {
	ProduceToKafka(req)
}

func ProduceToKafka(req *http.Request)  {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": c.KafkaCfg.Brokers})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := mux.Vars(req)["topic"]
	for _, word := range []string{"test message"} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func createHandlers() http.Handler{
	r := mux.NewRouter()
	for k := range c.ApiCfg {
		r.HandleFunc(fmt.Sprintf("/api/%s/{topic}",k), SendToKafka)
	}
	return r
}

func main() {
	fmt.Println("Go REST API for Kafka")
	configFile := flag.String("c", "config.yml", "go-kafka-rest configFile file")

	c.ReadConfig(*configFile)
	http.Handle("/", createHandlers())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
