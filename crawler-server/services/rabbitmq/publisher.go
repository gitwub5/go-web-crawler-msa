package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// PublishMessage는 지정된 큐에 메시지를 발행합니다.
func PublishMessage(queueName string, message string, ttl int) error {
	if channel == nil {
		return fmt.Errorf("RabbitMQ 채널이 초기화되지 않았습니다")
	}

	// 메시지 발행
	err := channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Expiration:  fmt.Sprintf("%d", ttl), // 메시지 TTL (밀리초 단위)
		},
	)
	if err != nil {
		log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		return err
	}

	log.Println("RabbitMQ 메시지 발행 성공:", message)
	return nil
}
