package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gitwub5/go-notification-web-server/config"
	"github.com/gitwub5/go-notification-web-server/rabbitmq"
)

// 공지사항 데이터를 관리하는 전역 변수
var notices []string
var mutex = &sync.Mutex{}

func main() {
	// 환경 변수 로드
	config.LoadEnv()

	// RabbitMQ 큐 목록
	queues := config.GetQueueNames()

	// RabbitMQ 구독 시작
	go rabbitmq.StartRabbitMQSubscribers(queues)

	// 새로운 메시지 처리
	go func() {
		for msg := range rabbitmq.NoticeChannel {
			mutex.Lock()
			notices = append(notices, msg) // 새 공지사항 추가
			log.Printf("새 공지사항 추가: %s", msg)
			mutex.Unlock()
		}
	}()

	// Gin 서버 설정
	r := gin.Default()
	r.LoadHTMLGlob("templates/*") // HTML 템플릿 경로 설정

	// 메인 페이지 렌더링
	r.GET("/", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		c.HTML(http.StatusOK, "index.html", gin.H{"notices": notices})
	})

	// JSON API로 공지사항 데이터 제공
	r.GET("/api/notices", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()

		// 공지사항을 분리하여 JSON 응답 생성
		c.JSON(http.StatusOK, gin.H{
			"cseNotices": notices, // CSE 공지사항 리스트
			"swNotices":  notices, // SW 공지사항 리스트
		})
	})

	// 서버 실행
	port := config.GetPort()
	log.Printf("서버 실행 중: http://localhost:%s", port)
	r.Run(":" + port)
}
