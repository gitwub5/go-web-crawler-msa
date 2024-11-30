package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv는 환경 변수를 로드합니다.
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env 파일을 찾을 수 없습니다. 시스템 환경 변수를 사용합니다.")
	}
}

// GetRabbitMQURL는 RabbitMQ URL을 반환합니다.
func GetRabbitMQURL() string {
	return os.Getenv("RABBITMQ_URL")
}

// GetQueueNames는 RabbitMQ 큐 이름 목록을 반환합니다.
func GetQueueNames() []string {
	queues := os.Getenv("QUEUE_NAMES")
	return strings.Split(queues, ",")
}

// GetPort는 서버 실행 포트를 반환합니다.
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 기본 포트
	}
	return port
}
