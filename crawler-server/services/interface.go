package services

import "github.com/JinHyeokOh01/go-crwl-server/models"

// GetCSECrawlingData는 컴퓨터공학과 공지사항을 크롤링하여 데이터를 반환합니다.
func GetCSECrawlingData() ([]models.Notice, error) {
	return CrwlCSENotices(cseURL)
}

// GetSWCrawlingData는 소프트웨어중심대학사업단 공지사항을 크롤링하여 데이터를 반환합니다.
func GetSWCrawlingData() ([]models.Notice, error) {
	swUrl := "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"
	return CrwlSWNotices(swUrl)
}
