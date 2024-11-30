package rabbitmq

import (
	"log"
	"sync"

	"github.com/streadway/amqp"
)

// NoticeChannel은 새 공지사항 메시지를 전달하기 위한 채널입니다.
var NoticeChannel = make(chan string)

// StartRabbitMQSubscriber는 RabbitMQ의 여러 큐를 구독합니다.
func StartRabbitMQSubscribers(queues []string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("RabbitMQ 연결 실패: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ 채널 생성 실패: %v", err)
	}
	defer ch.Close()

	var wg sync.WaitGroup

	// 각 큐에 대해 구독 설정
	for _, queue := range queues {
		wg.Add(1)
		go func(queueName string) {
			defer wg.Done()

			_, err := ch.QueueDeclare(
				queueName, // 큐 이름
				true,      // durable
				false,     // delete when unused
				false,     // exclusive
				false,     // no-wait
				nil,       // arguments
			)
			if err != nil {
				log.Printf("RabbitMQ 큐 선언 실패: %v", err)
				return
			}

			// 메시지 구독
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
				log.Printf("RabbitMQ 메시지 구독 실패 (%s): %v", queueName, err)
				return
			}

			log.Printf("%s 큐 구독 시작...", queueName)

			for msg := range msgs {
				log.Printf("[%s 큐] 새 메시지 수신: %s", queueName, msg.Body)
				NoticeChannel <- string(msg.Body)
			}
		}(queue)
	}

	wg.Wait()
}
