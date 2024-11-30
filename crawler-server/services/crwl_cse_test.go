package services

import (
	"log"
	"testing"
)

func TestCrwlCSENotices(t *testing.T) {
	// 테스트용 URL (CSE 공지사항 URL)
	url := "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

	// 크롤링 함수 호출
	notices, err := CrwlCSENotices(url)
	if err != nil {
		t.Fatalf("크롤링 중 오류 발생: %v", err)
	}

	// 결과 출력
	log.Printf("크롤링된 공지사항 개수: %d", len(notices))
	for _, notice := range notices {
		log.Printf("공지사항: %+v", notice)
	}

	// 결과 검증
	if len(notices) == 0 {
		t.Errorf("공지사항이 크롤링되지 않았습니다")
	}
}
