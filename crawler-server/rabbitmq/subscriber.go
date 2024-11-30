package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// SubscribeMessages는 지정된 큐에서 메시지를 구독합니다.
func SubscribeMessages(queueName string) (<-chan amqp.Delivery, error) {
	// RabbitMQ 연결
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 채널 생성
	ch, err := GetChannel(conn)
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	// 큐 선언
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Println("RabbitMQ 큐 선언 실패:", err)
		return nil, err
	}

	// 메시지 소비
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Println("RabbitMQ 메시지 구독 실패:", err)
		return nil, err
	}

	log.Println("RabbitMQ 메시지 구독 시작:", queueName)
	return msgs, nil
}
