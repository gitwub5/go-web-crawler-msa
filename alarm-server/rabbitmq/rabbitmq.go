package main

import (
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ 설정
const (
	rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/"
	cseQueue    = "cse-notices"
	swQueue     = "sw-notices"
)

func main() {
	// RabbitMQ 연결
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
	}
	defer conn.Close()

	// 채널 생성
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ 채널 생성 실패: %v", err)
	}
	defer ch.Close()

	// 큐 소비자 생성 (CSE 공지사항)
	go consumeQueue(ch, cseQueue)

	// 큐 소비자 생성 (SW 공지사항)
	go consumeQueue(ch, swQueue)

	// 종료되지 않도록 무한 대기
	select {}
}

// 큐에서 메시지를 소비하는 함수
func consumeQueue(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // 큐 이름
		"",        // 소비자 이름 (빈 문자열: RabbitMQ가 자동 생성)
		true,      // 자동 ACK
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("큐 소비자 설정 실패 (%s): %v", queueName, err)
	}

	log.Printf("큐 소비 시작: %s", queueName)

	// 메시지 수신
	for msg := range msgs {
		log.Printf("[%s 큐] 메시지 수신: %s", queueName, msg.Body)
		// 메시지 처리 로직 추가
		handleMessage(queueName, string(msg.Body))
	}
}

// 메시지를 처리하는 함수
func handleMessage(queueName, message string) {
	log.Printf("[%s 큐] 메시지 처리: %s", queueName, message)
	// 알림 처리 로직 (예: Redis 캐싱, MySQL 저장 등) 추가
}
