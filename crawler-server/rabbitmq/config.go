package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func GetConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("RabbitMQ 연결 실패:", err)
		return nil, err
	}
	return conn, nil
}

func GetChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("RabbitMQ 채널 생성 실패:", err)
		return nil, err
	}
	return ch, nil
}
