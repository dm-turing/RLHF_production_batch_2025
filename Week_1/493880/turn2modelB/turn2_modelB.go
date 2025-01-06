package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

// Subject represents an object with observers
type Subject interface {
	RegisterObserver(Observer)
	DeregisterObserver(Observer)
	NotifyObservers()
}

// ConcreteSubject is a specific subject that manages observers
type ConcreteSubject struct {
	observers []Observer
	mutex     sync.Mutex
}

func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{}
}
func (cs *ConcreteSubject) RegisterObserver(o Observer) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.observers = append(cs.observers, o)
}
func (cs *ConcreteSubject) DeregisterObserver(o Observer) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for i, observer := range cs.observers {
		if observer == o {
			cs.observers = append(cs.observers[:i], cs.observers[i+1:]...)
			break
		}
	}
}
func (cs *ConcreteSubject) NotifyObservers() {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	for _, observer := range cs.observers {
		observer.Update()
	}
}

// Observer represents an object that is observed
type Observer interface {
	Update()
}

// ConcreteObserver is a specific observer
type ConcreteObserver struct {
	name     string
	producer sarama.AsyncProducer
	topic    string
}

func NewConcreteObserver(name string, brokerList []string, topic string) *ConcreteObserver {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Retry.Max = 10
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}
	return &ConcreteObserver{name, producer, topic}
}
func (co *ConcreteObserver) Update() {
	// Publish a notification to the broker
	co.producer.Input() <- &sarama.ProducerMessage{
		Topic: co.topic,
		Value: sarama.StringEncoder("Notification from " + co.name),
	}
	// Receive producer errors from the channel
	go func() {
		for err := range co.producer.Errors() {
			fmt.Println("Producer Error:", err)
		}
	}()
}

// KafkaProducerSubject is a subject that uses Kafka as a message broker
type KafkaProducerSubject struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewKafkaProducerSubject(brokerList []string, topic string) *KafkaProducerSubject {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Retry.Max = 10
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		panic(err)
	}
	return &KafkaProducerSubject{producer, topic}
}
func (kps *KafkaProducerSubject) NotifyObservers() {
	// Publish a notification to the broker
	kps.producer.Input() <- &sarama.ProducerMessage{
		Topic: kps.topic,
		Value: sarama.StringEncoder("Notification from KafkaProducerSubject"),
	}
	// Receive producer errors from the channel
	go func() {
		for err := range kps.producer.Errors() {
			fmt.Println("Producer Error:", err)
		}
	}()
}
func main() {
	// Kafka broker list
	brokerList := []string{"localhost:9092"}
	// Topic name
	topic := "notifications"
	subject := NewKafkaProducerSubject(brokerList, topic)
	observer1 := NewConcreteObserver("Observer 1", brokerList, topic)
	observer2 := NewConcreteObserver("Observer 2", brokerList, topic)
	subject.RegisterObserver(observer1)
	subject.RegisterObserver(observer2)
	subject.NotifyObservers()
	// Wait for some time before exiting
	time.Sleep(5 * time.Second)
}
