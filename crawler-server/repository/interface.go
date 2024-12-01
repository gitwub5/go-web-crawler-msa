package repository

import "github.com/JinHyeokOh01/go-crwl-server/models"

// NoticeRepository 인터페이스 정의
type NoticeRepository interface {
	CreateBatchNotices(tableName string, notices []models.Notice) error
	CreateBatchCSE(notices []models.Notice) error
	CreateBatchSW(notices []models.Notice) error
	DeleteBatchNotices(tableName string, notices []models.Notice) error
	DeleteBatchCSE(notices []models.Notice) error
	DeleteBatchSW(notices []models.Notice) error
	GetAllNotices(tableName string) ([]models.Notice, error)
	GetAllCSENotices() ([]models.Notice, error)
	GetAllSWNotices() ([]models.Notice, error)
	DeleteAllNotices(tableName string) error
	DeleteAllCSE() error
	DeleteAllSW() error
	GetLatestNotice(tableName string) (models.Notice, error)
}
