package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// PublishMessage는 지정된 큐에 메시지를 발행합니다.
func PublishMessage(queueName string, message string) error {
	// RabbitMQ 연결
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	// 채널 생성
	ch, err := GetChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	// 큐 선언
	_, err = ch.QueueDeclare(
		queueName, // 큐 이름
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("RabbitMQ 큐 선언 실패:", err)
		return err
	}

	// 메시지 발행
	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println("RabbitMQ 메시지 발행 실패:", err)
		return err
	}

	log.Println("RabbitMQ 메시지 발행 성공:", message)
	return nil
}
