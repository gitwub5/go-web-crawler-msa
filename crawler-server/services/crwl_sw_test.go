package services

import (
	"log"
	"testing"
)

func TestCrwlSWNotices(t *testing.T) {
	// 테스트 대상 URL
	url := "https://swedu.khu.ac.kr/bbs/board.php?bo_table=07_01"

	// 크롤링 함수 호출
	notices, err := CrwlSWNotices(url)
	if err != nil {
		t.Fatalf("SW 공지사항 크롤링 중 오류 발생: %v", err)
	}

	// 결과 출력
	log.Printf("크롤링된 SW 공지사항 개수: %d", len(notices))
	for _, notice := range notices {
		log.Printf("공지사항: %+v", notice)
	}

	// 결과 검증
	if len(notices) == 0 {
		t.Errorf("크롤링된 SW 공지사항이 없습니다")
	}
}
