package controllers

import (
	"net/http"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/services"
	"github.com/gin-gonic/gin"
)

// GetCSENoticesHandler는 CSE 공지사항 크롤링 요청을 처리합니다.
func GetCSENoticesHandler(c *gin.Context) {
	// 서비스 계층 호출
	crawledNotices, err := services.GetCSECrawlingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "컴퓨터공학과 공지사항 크롤링 실패",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	// 성공 시 크롤링 결과 반환
	c.JSON(http.StatusOK, models.APIResponse{
		Message: "컴퓨터공학과 공지사항 크롤링 완료",
		Data:    crawledNotices,
		Error:   "",
	})
}

// GetSWNoticesHandler는 SW 공지사항 크롤링 요청을 처리합니다.
func GetSWNoticesHandler(c *gin.Context) {
	// 서비스 계층 호출
	crawledNotices, err := services.GetSWCrawlingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "소프트웨어중심대학사업단 공지사항 크롤링 실패",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	// 성공 시 크롤링 결과 반환
	c.JSON(http.StatusOK, models.APIResponse{
		Message: "소프트웨어중심대학사업단 공지사항 크롤링 완료",
		Data:    crawledNotices,
		Error:   "",
	})
}
