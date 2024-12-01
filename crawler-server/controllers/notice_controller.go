package controllers

import (
	"net/http"

	"github.com/JinHyeokOh01/go-crwl-server/services"
	"github.com/gin-gonic/gin"
)

type NoticeController struct {
	service services.NoticeService // 인터페이스로 변경
}

func NewNoticeController(service services.NoticeService) *NoticeController {
	return &NoticeController{
		service: service,
	}
}

// 유효한 테이블 이름을 정의
var validTables = map[string]bool{
	"cse_notices": true,
	"sw_notices":  true,
}

// validateTableName: 테이블 이름 검증
func validateTableName(tableName string) bool {
	return validTables[tableName]
}

// GetNotices: DB에서 공지사항 조회
func (nc *NoticeController) GetNotices(c *gin.Context) {
	tableName := c.Param("tableName") // URL 파라미터로 테이블 이름 전달

	// 테이블 이름 검증
	if !validateTableName(tableName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 테이블 이름입니다",
		})
		return
	}

	notices, err := nc.service.GetAllNotices(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": tableName + " 공지사항 조회 성공",
		"data":    notices,
	})
}

// DeleteAllNotices: DB의 모든 공지사항 삭제
func (nc *NoticeController) DeleteAllNotices(c *gin.Context) {
	tableName := c.Param("tableName") // URL 파라미터로 테이블 이름 전달

	// 테이블 이름 검증
	if !validateTableName(tableName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "유효하지 않은 테이블 이름입니다",
		})
		return
	}

	if err := nc.service.DeleteAllNotices(tableName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tableName + " 공지사항 삭제 실패: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": tableName + " 공지사항이 모두 삭제되었습니다",
	})
}
