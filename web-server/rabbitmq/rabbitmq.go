package rabbitmq

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

var (
	connection    *amqp.Connection
	channel       *amqp.Channel
	NoticeChannel = make(chan string) // 메시지를 전달하는 채널
	once          sync.Once           // 연결 초기화를 위한 싱글톤
)

// InitializeRabbitMQ initializes RabbitMQ connection and channel
func InitializeRabbitMQ(url string) error {
	var err error
	once.Do(func() {
		connection, err = GetConnection(url)
		if err != nil {
			log.Fatalf("RabbitMQ 연결 실패: %v", err)
			return
		}

		channel, err = GetChannel(connection)
		if err != nil {
			log.Fatalf("RabbitMQ 채널 생성 실패: %v", err)
			return
		}

		log.Println("RabbitMQ 초기화 성공")
	})
	return err
}

// GetConnection creates a new RabbitMQ connection
func GetConnection(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("RabbitMQ 연결 실패: %v", err)
		return nil, err
	}
	log.Println("RabbitMQ 연결 성공")
	return conn, nil
}

// GetChannel creates a new channel for RabbitMQ
func GetChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("RabbitMQ 채널 생성 실패: %v", err)
		return nil, err
	}
	log.Println("RabbitMQ 채널 생성 성공")
	return ch, nil
}

// StartRabbitMQSubscribers subscribes to RabbitMQ queues
func StartRabbitMQSubscribers(queues []string) {
	for _, queue := range queues {
		go subscribeToQueue(queue)
	}
}

func subscribeToQueue(queue string) {
	msgs, err := channel.Consume(
		queue, // 큐 이름
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("RabbitMQ 소비 실패 (%s): %v", queue, err)
	}

	for msg := range msgs {
		NoticeChannel <- string(msg.Body)

		// 수동 Ack 처리
		if err := msg.Ack(false); err != nil {
			log.Printf("Ack 처리 실패 (%s): %v", queue, err)
		} else {
			log.Printf("Ack 처리 성공 (%s)", queue)
		}
	}
}

// CloseRabbitMQ closes the RabbitMQ connection and channel
func CloseRabbitMQ() {
	if channel != nil {
		channel.Close()
	}
	if connection != nil {
		connection.Close()
	}
}
