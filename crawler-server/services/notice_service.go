package services

import "github.com/JinHyeokOh01/go-crwl-server/models"

// noticeService 구조체 정의
type noticeService struct {
	repo NoticeRepository
}

// NewNoticeService는 NoticeService 구현체를 반환합니다.
func NewNoticeService(repo NoticeRepository) NoticeService {
	return &noticeService{
		repo: repo,
	}
}

// GetAllNotices는 DB에서 공지사항을 조회합니다.
func (s *noticeService) GetAllNotices(tableName string) ([]models.Notice, error) {
	return s.repo.GetAllNotices(tableName)
}

// CreateBatchNotices는 공지사항을 DB에 저장합니다.
func (s *noticeService) CreateBatchNotices(tableName string, notices []models.Notice) error {
	return s.repo.CreateBatchNotices(tableName, notices)
}

// DeleteBatchNotices는 특정 공지사항을 DB에서 삭제합니다.
func (s *noticeService) DeleteBatchNotices(tableName string, notices []models.Notice) error {
	return s.repo.DeleteBatchNotices(tableName, notices)
}

// DeleteAllNotices는 특정 테이블의 모든 공지사항을 삭제합니다.
func (s *noticeService) DeleteAllNotices(tableName string) error {
	return s.repo.DeleteAllNotices(tableName)
}
