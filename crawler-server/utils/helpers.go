package utils

import "github.com/JinHyeokOh01/go-crwl-server/models"

func RemoveDuplicateNotices(notices []models.Notice, existing []models.Notice) []models.Notice {
	noticeMap := make(map[string]bool)
	for _, n := range existing {
		noticeMap[n.Number] = true // 기존 공지사항의 번호를 맵에 저장
	}

	uniqueNotices := []models.Notice{}
	for _, n := range notices {
		if !noticeMap[n.Number] { // 기존에 없는 공지사항만 추가
			uniqueNotices = append(uniqueNotices, n)
			noticeMap[n.Number] = true
		}
	}
	return uniqueNotices
}
