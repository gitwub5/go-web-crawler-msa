package services

import (
	"log"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/repository"
	"github.com/JinHyeokOh01/go-crwl-server/services/crwl"
	"github.com/JinHyeokOh01/go-crwl-server/services/rabbitmq"
	"github.com/JinHyeokOh01/go-crwl-server/utils"
)

// crawlingService 구조체
type crawlingService struct {
	repo repository.NoticeRepository
}

// NewCrawlingService는 CrawlingService 구현체를 생성합니다.
func NewCrawlingService(repo repository.NoticeRepository) CrawlingService {
	return &crawlingService{
		repo: repo,
	}
}

// HandleCSECrawling은 컴퓨터공학과 공지사항 크롤링, 데이터베이스 저장, RabbitMQ 발행을 수행합니다.
func (s *crawlingService) HandleCSECrawling() ([]models.Notice, error) {
	const cseURL = "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

	// 크롤링 수행
	crawledNotices, err := crwl.CrwlCSENotices(cseURL)
	if err != nil {
		log.Printf("CSE 공지사항 크롤링 실패: %v", err)
		return nil, err
	}

	// 최신 공지사항 가져오기
	latestNotice, err := s.repo.GetLatestNotice("cse_notices")
	if err != nil {
		log.Printf("CSE 최근 공지사항 조회 실패: %v", err)
		return nil, err
	}

	// 최신 데이터 필터링
	newNotices := filterNewNoticesAfter(latestNotice, crawledNotices)

	// 새로운 공지사항이 있을 경우 처리
	if len(newNotices) > 0 {
		log.Printf("새로운 CSE 공지사항 %d개 발견", len(newNotices))

		// DB에 저장
		if err := s.repo.CreateBatchNotices("cse_notices", newNotices); err != nil {
			log.Printf("CSE 공지사항 저장 실패: %v", err)
			return nil, err
		}

		// RabbitMQ로 발행
		for _, notice := range newNotices {
			message := utils.FormatNoticeMessage(notice)
			if err := rabbitmq.PublishMessage("cse-notices", message, 10000); err != nil { // TTL: 10초
				log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
			} else {
				log.Printf("RabbitMQ 메시지 발행 성공: %s", message)
			}
		}
	} else {
		// 새로운 공지사항이 없을 경우 메시지 발행
		message := utils.FormatNoNewNoticesMessage("cse-notices")
		if err := rabbitmq.PublishMessage("cse-notices", message, 10000); err != nil {
			log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		} else {
			log.Printf("RabbitMQ 메시지 발행 성공: %s", message)
		}
		log.Println("새로운 CSE 공지사항이 없습니다.")
	}

	return crawledNotices, nil
}

// HandleSWCrawling은 소프트웨어중심대학사업단 공지사항 크롤링, 데이터베이스 저장, RabbitMQ 발행을 수행합니다.
func (s *crawlingService) HandleSWCrawling() ([]models.Notice, error) {
	const swURL = "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"

	// 크롤링 수행
	crawledNotices, err := crwl.CrwlSWNotices(swURL)
	if err != nil {
		log.Printf("SW 공지사항 크롤링 실패: %v", err)
		return nil, err
	}

	// 최신 공지사항 가져오기
	latestNotice, err := s.repo.GetLatestNotice("sw_notices")
	if err != nil {
		log.Printf("SW 최근 공지사항 조회 실패: %v", err)
		return nil, err
	}

	// 최신 데이터 필터링
	newNotices := filterNewNoticesAfter(latestNotice, crawledNotices)

	// 새로운 공지사항이 있을 경우 처리
	if len(newNotices) > 0 {
		log.Printf("새로운 SW 공지사항 %d개 발견", len(newNotices))

		// DB에 저장
		if err := s.repo.CreateBatchNotices("sw_notices", newNotices); err != nil {
			log.Printf("SW 공지사항 저장 실패: %v", err)
			return nil, err
		}

		// RabbitMQ로 발행
		for _, notice := range newNotices {
			message := utils.FormatNoticeMessage(notice)
			if err := rabbitmq.PublishMessage("sw-notices", message, 10000); err != nil {
				log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
			} else {
				log.Printf("RabbitMQ 메시지 발행 성공: %s", message)
			}
		}
	} else {
		// 새로운 공지사항이 없을 경우 메시지 발행
		message := utils.FormatNoNewNoticesMessage("sw-notices")
		if err := rabbitmq.PublishMessage("sw-notices", message, 10000); err != nil {
			log.Printf("RabbitMQ 메시지 발행 실패: %v", err)
		} else {
			log.Printf("RabbitMQ 메시지 발행 성공: %s", message)
		}
		log.Println("새로운 SW 공지사항이 없습니다.")
	}

	return crawledNotices, nil
}

// filterNewNoticesAfter는 최신 공지사항 이후의 데이터를 필터링합니다.
func filterNewNoticesAfter(latest models.Notice, crawled []models.Notice) []models.Notice {
	var newNotices []models.Notice

	for _, notice := range crawled {
		if notice.Date > latest.Date || (notice.Date == latest.Date && notice.Number > latest.Number) {
			newNotices = append(newNotices, notice)
		}
	}

	return newNotices
}
