package controllers

import (
	"net/http"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/JinHyeokOh01/go-crwl-server/services"
	"github.com/gin-gonic/gin"
)

type CrwlController struct {
	crawlingService services.CrawlingService
}

func NewCrwlController(crawlingService services.CrawlingService) *CrwlController {
	return &CrwlController{crawlingService: crawlingService}
}

// HandleCSECrawling 트리거를 처리합니다.
func (cc *CrwlController) HandleCSECrawling(c *gin.Context) {
	// 크롤링 서비스 호출
	crawledNotices, err := cc.crawlingService.HandleCSECrawling()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "컴퓨터공학과 공지사항 크롤링 실패",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "컴퓨터공학과 공지사항 크롤링 완료",
		Data:    crawledNotices,
		Error:   "",
	})
}

// HandleSWCrawling 트리거를 처리합니다.
func (cc *CrwlController) HandleSWCrawling(c *gin.Context) {
	// 크롤링 서비스 호출
	crawledNotices, err := cc.crawlingService.HandleSWCrawling()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Message: "소프트웨어중심대학사업단 공지사항 크롤링 실패",
			Data:    nil,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Message: "소프트웨어중심대학사업단 공지사항 크롤링 완료",
		Data:    crawledNotices,
		Error:   "",
	})
}
