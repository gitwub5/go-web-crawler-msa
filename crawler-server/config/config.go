package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func LoadEnv() {
	// .env 파일 로드
	if err := godotenv.Load(".env"); err != nil {
		log.Println(".env 파일을 로드하지 못했습니다. 환경 변수를 직접 사용합니다.")
	}
}

// GetDBConfig는 데이터베이스 설정을 반환합니다.
func GetDBConfig() DBConfig {
	port, err := strconv.Atoi(getEnv("MYSQL_PORT", "3306"))
	if err != nil {
		log.Fatalf("MySQL 포트 변환 실패: %v", err)
	}

	return DBConfig{
		Host:     getEnv("MYSQL_HOST", "localhost"),
		Port:     port,
		User:     getEnv("MYSQL_USER", "root"),
		Password: getEnv("MYSQL_PASSWORD", "password"),
		Name:     getEnv("MYSQL_DATABASE", "crwl_db"),
	}
}

// GetRabbitMQURL는 RabbitMQ 연결 URL을 반환합니다.
func GetRabbitMQURL() string {
	url := getEnv("RABBITMQ_URL", "")
	if url == "" {
		log.Fatalf("RabbitMQ URL이 설정되지 않았습니다.")
	}
	return url
}

// GetQueueNames는 RabbitMQ 큐 이름을 쉼표로 구분한 목록으로 반환합니다.
func GetQueueNames() []string {
	queues := getEnv("QUEUE_NAMES", "")
	if queues == "" {
		log.Fatalf("RabbitMQ 큐 이름이 설정되지 않았습니다.")
	}
	return strings.Split(queues, ",")
}

// GetPort는 크롤링 서버 실행 포트를 반환합니다.
func GetPort() string {
	port := getEnv("CRAWLER_SERVER_PORT", "8080") // 기본 포트
	if _, err := strconv.Atoi(port); err != nil {
		log.Fatalf("유효하지 않은 PORT 값: %v", port)
	}
	return port
}

// getEnv는 환경 변수의 값을 가져오거나, 없으면 기본값을 반환합니다.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
