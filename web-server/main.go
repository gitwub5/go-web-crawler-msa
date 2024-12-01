package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gitwub5/go-notification-web-server/config"
	"github.com/gitwub5/go-notification-web-server/rabbitmq"
)

// Notice 구조체
type Notice struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Link   string `json:"link"`
}

// NoticesResponse 구조체
type NoticesResponse struct {
	Data    []Notice `json:"data"`
	Message string   `json:"message"`
}

// 공지사항 데이터를 관리하는 전역 변수
var (
	cseNotices        []Notice
	swNotices         []Notice
	notificationQueue []string
	notificationMutex = &sync.Mutex{}
	noticesMutex      = &sync.Mutex{}
)

func main() {
	// 환경 변수 로드
	config.LoadEnv()

	// RabbitMQ 초기화
	rabbitMQURL := config.GetRabbitMQURL()
	if err := rabbitmq.InitializeRabbitMQ(rabbitMQURL); err != nil {
		log.Fatalf("RabbitMQ 초기화 실패: %v", err)
	}
	defer rabbitmq.CloseRabbitMQ()

	// RabbitMQ 구독 시작
	queues := []string{"cse-notices", "sw-notices"}
	rabbitmq.StartRabbitMQSubscribers(queues)

	// RabbitMQ 메시지 처리
	go handleRabbitMQMessages()

	// Gin 서버 설정
	r := gin.Default()
	r.LoadHTMLGlob("templates/*") // HTML 템플릿 경로 설정

	// 메인 페이지 렌더링
	r.GET("/", func(c *gin.Context) {
		noticesMutex.Lock()
		defer noticesMutex.Unlock()
		notificationMutex.Lock()
		defer notificationMutex.Unlock()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"cseNotices":          cseNotices,
			"swNotices":           swNotices,
			"latestNotifications": notificationQueue,
		})
	})

	// JSON API로 공지사항 데이터 제공
	r.GET("/api/notices", func(c *gin.Context) {
		noticesMutex.Lock()
		defer noticesMutex.Unlock()
		notificationMutex.Lock()
		defer notificationMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"cseNotices":          cseNotices,
			"swNotices":           swNotices,
			"latestNotifications": notificationQueue,
		})

		// 알림 초기화
		notificationQueue = []string{}
	})

	// 초기 데이터 로드
	err := fetchNoticesFromCrawler()
	if err != nil {
		log.Printf("초기 데이터 로드 중 오류 발생: %v", err)
	}

	// 서버 실행
	port := config.GetPort()
	log.Printf("웹 서버 실행 중: http://localhost:%s", port)
	r.Run(":" + port)
}

// handleRabbitMQMessages processes messages from RabbitMQ
func handleRabbitMQMessages() {
	for msg := range rabbitmq.NoticeChannel {
		log.Printf("RabbitMQ 메시지 수신: %s", msg)

		// 메시지를 알림 큐에 추가
		notificationMutex.Lock()
		notificationQueue = append(notificationQueue, msg)
		notificationMutex.Unlock()
	}
}

// fetchNoticesFromCrawler fetches notices from the crawler server
func fetchNoticesFromCrawler() error {
	crawlerURL := config.GetCrawlerServerURL()
	if crawlerURL == "" {
		log.Println("크롤러 서버 URL이 설정되지 않았습니다.")
		return fmt.Errorf("크롤러 서버 URL이 설정되지 않았습니다")
	}

	cseURL := crawlerURL + "/notices/cse_notices"
	swURL := crawlerURL + "/notices/sw_notices"

	cseData, err := fetchNotices(cseURL)
	if err != nil {
		return err
	}

	swData, err := fetchNotices(swURL)
	if err != nil {
		return err
	}

	noticesMutex.Lock()
	defer noticesMutex.Unlock()
	cseNotices = cseData
	swNotices = swData

	log.Println("크롤러 서버에서 데이터 가져오기 완료")
	return nil
}

// fetchNotices fetches notices from the given URL
func fetchNotices(url string) ([]Notice, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("공지사항을 가져오는 중 오류 발생: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("공지사항 응답 읽기 오류: %v", err)
		return nil, err
	}

	var response NoticesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("공지사항 JSON 파싱 오류: %v", err)
		return nil, err
	}

	return response.Data, nil
}
