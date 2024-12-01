package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
)

// InitializeRabbitMQ는 RabbitMQ 연결 및 채널을 초기화합니다.
func InitializeRabbitMQ(url string) error {
	var err error

	// RabbitMQ 연결 생성
	connection, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
		return err
	}
	log.Println("RabbitMQ 연결 성공")

	// RabbitMQ 채널 생성
	channel, err = connection.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ 채널 생성 실패: %v", err)
		connection.Close()
		return err
	}
	log.Println("RabbitMQ 채널 생성 성공")

	// 필요한 큐 생성
	queueNames := []string{"cse-notices", "sw-notices"}
	for _, queueName := range queueNames {
		if err := CreateQueue(queueName); err != nil {
			log.Fatalf("RabbitMQ 큐 생성 실패 (%s): %v", queueName, err)
			CloseRabbitMQ()
			return err
		}
	}

	return nil
}

// CreateQueue는 주어진 이름의 큐를 생성합니다.
func CreateQueue(queueName string) error {
	_, err := channel.QueueDeclare(
		queueName, // 큐 이름
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	log.Printf("RabbitMQ 큐 생성 성공: %s", queueName)
	return nil
}

// CloseRabbitMQ는 RabbitMQ 연결과 채널을 닫습니다.
func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}
