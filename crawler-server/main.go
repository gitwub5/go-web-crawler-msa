package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/config"
	"github.com/JinHyeokOh01/go-crwl-server/controllers"
	"github.com/JinHyeokOh01/go-crwl-server/repository"
	"github.com/JinHyeokOh01/go-crwl-server/services"
	"github.com/JinHyeokOh01/go-crwl-server/services/rabbitmq"
	"github.com/JinHyeokOh01/go-crwl-server/store"

	"github.com/gin-gonic/gin"
)

func startPeriodicCrawling(crawlingService services.CrawlingService) {
	// ticker := time.NewTicker(1 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("주기적 크롤링 시작")

		// CSE 공지사항 크롤링 및 저장
		_, err := crawlingService.HandleCSECrawling()
		if err != nil {
			log.Printf("CSE 공지사항 크롤링 중 오류 발생: %v", err)
		}

		// SW 공지사항 크롤링 및 저장
		_, err = crawlingService.HandleSWCrawling()
		if err != nil {
			log.Printf("SW 공지사항 크롤링 중 오류 발생: %v", err)
		}

		log.Println("주기적 크롤링 완료")
	}
}

func main() {
	// 환경 변수 로드
	config.LoadEnv()

	// 데이터베이스 초기화
	dbConfig := config.GetDBConfig()
	if err := store.Initialize(dbConfig); err != nil {
		log.Fatalf("데이터베이스 초기화 실패: %v", err)
	}
	defer store.Close()

	// RabbitMQ 초기화
	rabbitMQURL := config.GetRabbitMQURL()
	if err := rabbitmq.InitializeRabbitMQ(rabbitMQURL); err != nil {
		log.Fatalf("RabbitMQ 초기화 실패: %v", err)
	}
	defer rabbitmq.CloseRabbitMQ()

	// 애플리케이션 종료 시그널 처리
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 레포지토리 및 서비스 생성
	repo := repository.NewNoticeRepository()
	noticeService := services.NewNoticeService(repo)
	crawlingService := services.NewCrawlingService(repo)

	// Controller 생성
	noticeController := controllers.NewNoticeController(noticeService)
	crwlController := controllers.NewCrwlController(crawlingService)

	// Gin 서버 설정
	r := gin.Default()

	// API 라우팅 설정
	r.GET("/crawling/cse", crwlController.HandleCSECrawling)
	r.GET("/crawling/sw", crwlController.HandleSWCrawling)
	r.GET("/notices/:tableName", noticeController.GetNotices)
	r.DELETE("/notices/:tableName", noticeController.DeleteAllNotices)

	// 주기적 크롤링 시작
	go startPeriodicCrawling(crawlingService)

	// 서버 실행
	go func() {
		log.Println("API 서버 실행 중...")
		if err := r.Run(":" + config.GetPort()); err != nil {
			log.Fatalf("서버 실행 중 오류 발생: %v", err)
		}
	}()

	// 종료 신호 처리
	<-sigs
	log.Println("애플리케이션 종료 중...")
}
