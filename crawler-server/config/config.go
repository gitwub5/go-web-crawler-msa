package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("환경변수 로드 실패: %v", err)
	}
}

// GetRabbitMQURL는 RabbitMQ 연결 URL을 반환합니다.
func GetRabbitMQURL() string {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatalf("RabbitMQ URL이 설정되지 않았습니다.")
	}
	return url
}

// GetQueueNames는 RabbitMQ 큐 이름을 쉼표로 구분한 목록으로 반환합니다.
func GetQueueNames() []string {
	queues := os.Getenv("QUEUE_NAMES")
	if queues == "" {
		log.Fatalf("RabbitMQ 큐 이름이 설정되지 않았습니다.")
	}
	return strings.Split(queues, ",")
}

// GetPort는 크롤링 서버 실행 포트를 반환합니다.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // 기본 포트
	}
	return port
}
