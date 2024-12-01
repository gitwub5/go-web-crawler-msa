package services

import "github.com/JinHyeokOh01/go-crwl-server/models"

// NoticeService 인터페이스: 공지사항 관련 비즈니스 로직을 정의합니다.
type NoticeService interface {
	GetAllNotices(tableName string) ([]models.Notice, error)
	CreateBatchNotices(tableName string, notices []models.Notice) error
	DeleteBatchNotices(tableName string, notices []models.Notice) error
	DeleteAllNotices(tableName string) error
}

// CrawlingService 인터페이스: 크롤링 관련 비즈니스 로직을 정의합니다.
type CrawlingService interface {
	HandleCSECrawling() ([]models.Notice, error) // CSE 공지사항 크롤링 및 처리
	HandleSWCrawling() ([]models.Notice, error)  // SW 공지사항 크롤링 및 처리
}

// NoticeRepository 인터페이스: 공지사항 관련 데이터 접근 계층을 정의합니다.
type NoticeRepository interface {
	GetAllNotices(tableName string) ([]models.Notice, error)
	CreateBatchNotices(tableName string, notices []models.Notice) error
	DeleteBatchNotices(tableName string, notices []models.Notice) error
	DeleteAllNotices(tableName string) error
	GetLatestNotice(tableName string) (models.Notice, error)
}
