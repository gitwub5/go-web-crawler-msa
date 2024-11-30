package main

import (
	"log"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/config"
	"github.com/JinHyeokOh01/go-crwl-server/controllers"
	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/rabbitmq"
	"github.com/JinHyeokOh01/go-crwl-server/services"
	"github.com/JinHyeokOh01/go-crwl-server/utils"

	"github.com/gin-gonic/gin"
)

func startPeriodicCrawling() {
	// 각 공지사항의 기존 데이터 저장을 위한 슬라이스
	existingCSENotices := []models.Notice{}
	existingSWNotices := []models.Notice{}

	var executeCrawling func()
	executeCrawling = func() {
		// CSE 공지사항 크롤링 및 발행
		updatedCSE := handleCrawlingAndPublish("cse-notices", services.GetCSECrawlingData, existingCSENotices)
		existingCSENotices = updatedCSE

		// SW 공지사항 크롤링 및 발행
		updatedSW := handleCrawlingAndPublish("sw-notices", services.GetSWCrawlingData, existingSWNotices)
		existingSWNotices = updatedSW

		// 다음 크롤링 예약
		time.AfterFunc(1*time.Minute, executeCrawling)
	}

	// 첫 번째 크롤링 시작
	executeCrawling()
}

// handleCrawlingAndPublish는 크롤링 데이터를 가져오고 RabbitMQ에 발행하며, 중복 제거된 공지사항을 반환합니다.
func handleCrawlingAndPublish(queueName string, crawlingFunc func() ([]models.Notice, error), existing []models.Notice) []models.Notice {
	notices, err := crawlingFunc()
	if err != nil {
		log.Printf("[%s] 크롤링 실패: %v", queueName, err)
		return existing
	}

	// 중복 제거
	uniqueNotices := utils.RemoveDuplicateNotices(notices, existing)

	// 새로운 공지사항이 있을 경우에만 로그 출력
	if len(uniqueNotices) > 0 {
		log.Printf("[%s] 새로운 공지사항 %d개 발견", queueName, len(uniqueNotices))
	} else {
		log.Printf("[%s] 새로운 공지사항이 없습니다", queueName)
	}

	// RabbitMQ 발행
	for _, notice := range uniqueNotices {
		message := utils.FormatNoticeMessage(notice)
		err := rabbitmq.PublishMessage(queueName, message)
		if err != nil {
			log.Printf("[%s] RabbitMQ 발행 실패: %v", queueName, err)
		} else {
			log.Printf("[%s] RabbitMQ 발행 성공: %s", queueName, message)
		}
	}

	// 기존 공지사항 목록을 최신화
	return append(existing, uniqueNotices...)
}

func main() {
	// 환경 변수 로드
	config.LoadEnv()

	port := config.GetPort() // 환경 변수에서 포트 가져오기
	log.Printf("크롤링 서버 실행 중: http://localhost:%s", port)

	// Gin 서버 설정
	r := gin.Default()

	// CSE 공지사항 엔드포인트 연결
	r.GET("/cse", controllers.GetCSENoticesHandler)
	// SW 공지사항 엔드포인트 연결
	r.GET("/sw", controllers.GetSWNoticesHandler)

	// 주기적 크롤링 시작
	go startPeriodicCrawling()

	// 서버 실행
	r.Run(":" + port)
}
