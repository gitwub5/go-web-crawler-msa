package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/store"
)

const (
	CSETableName = "cse_notices"
	SWTableName  = "sw_notices"
)

// SQLNoticeRepository 구현체

type noticeRepository struct {
	db *sql.DB
}

// NewSQLNoticeRepository는 SQLNoticeRepository의 인스턴스를 생성합니다.
func NewNoticeRepository() NoticeRepository {
	return &noticeRepository{
		db: store.DB, // store.DB는 초기화된 *sql.DB 객체여야 합니다.
	}
}

// CreateBatchNotices 공지사항 일괄 저장
func (r *noticeRepository) CreateBatchNotices(tableName string, notices []models.Notice) error {
	if len(notices) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
        INSERT INTO %s (number, title, date, link)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE
            title = VALUES(title),
            date = VALUES(date),
            link = VALUES(link)
    `, tableName)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, notice := range notices {
		_, err = stmt.Exec(notice.Number, notice.Title, notice.Date, notice.Link)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CreateBatchCSE는 CSE 공지사항을 저장합니다.
func (r *noticeRepository) CreateBatchCSE(notices []models.Notice) error {
	return r.CreateBatchNotices(CSETableName, notices)
}

// CreateBatchSW는 SW 공지사항을 저장합니다.
func (r *noticeRepository) CreateBatchSW(notices []models.Notice) error {
	return r.CreateBatchNotices(SWTableName, notices)
}

// DeleteBatchNotices 공지사항 일괄 삭제
func (r *noticeRepository) DeleteBatchNotices(tableName string, notices []models.Notice) error {
	if len(notices) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	placeholders := make([]string, len(notices))
	args := make([]interface{}, len(notices))

	for i, notice := range notices {
		placeholders[i] = "?"
		args[i] = notice.Number
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE number IN (%s)", tableName, strings.Join(placeholders, ","))
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteBatchCSE는 CSE 공지사항을 삭제합니다.
func (r *noticeRepository) DeleteBatchCSE(notices []models.Notice) error {
	return r.DeleteBatchNotices(CSETableName, notices)
}

// DeleteBatchSW는 SW 공지사항을 삭제합니다.
func (r *noticeRepository) DeleteBatchSW(notices []models.Notice) error {
	return r.DeleteBatchNotices(SWTableName, notices)
}

// GetAllNotices 공지사항 전체 조회
func (r *noticeRepository) GetAllNotices(tableName string) ([]models.Notice, error) {
	query := fmt.Sprintf(`
        SELECT number, title, date, link 
        FROM %s 
        ORDER BY date DESC, number DESC
    `, tableName)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notices []models.Notice
	for rows.Next() {
		var n models.Notice
		if err := rows.Scan(&n.Number, &n.Title, &n.Date, &n.Link); err != nil {
			return nil, err
		}
		notices = append(notices, n)
	}
	return notices, nil
}

// GetAllCSENotices는 CSE 공지사항을 조회합니다.
func (r *noticeRepository) GetAllCSENotices() ([]models.Notice, error) {
	return r.GetAllNotices(CSETableName)
}

// GetAllSWNotices는 SW 공지사항을 조회합니다.
func (r *noticeRepository) GetAllSWNotices() ([]models.Notice, error) {
	return r.GetAllNotices(SWTableName)
}

// DeleteAllNotices 공지사항 전체 삭제
func (r *noticeRepository) DeleteAllNotices(tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s", tableName)
	_, err := r.db.Exec(query)
	return err
}

// DeleteAllCSE는 CSE 공지사항을 삭제합니다.
func (r *noticeRepository) DeleteAllCSE() error {
	return r.DeleteAllNotices(CSETableName)
}

// DeleteAllSW는 SW 공지사항을 삭제합니다.
func (r *noticeRepository) DeleteAllSW() error {
	return r.DeleteAllNotices(SWTableName)
}

// GetLatestNotice는 가장 최신 공지사항을 조회합니다.
func (r *noticeRepository) GetLatestNotice(tableName string) (models.Notice, error) {
	var notice models.Notice

	query := `
        SELECT number, title, date, link 
        FROM ` + tableName + `
        ORDER BY date DESC, number DESC
        LIMIT 1
    `

	err := r.db.QueryRow(query).Scan(&notice.Number, &notice.Title, &notice.Date, &notice.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			// 테이블에 데이터가 없으면 빈 공지사항 반환
			return models.Notice{}, nil
		}
		return notice, err
	}

	return notice, nil
}
