package utils

import (
	"fmt"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/models"
)

// FormatNoticeMessage는 Notice 데이터를 문자열로 포맷합니다.
func FormatNoticeMessage(notice models.Notice) string {
	return notice.Number + "|" + notice.Title + "|" + notice.Date + "|" + notice.Link
}

// FormatNoNewNoticesMessage는 새로운 공지사항이 없을 때 발행할 메시지를 생성합니다.
func FormatNoNewNoticesMessage(noticeType string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] 새로운 공지사항이 없습니다. (업데이트 시간: %s)", noticeType, currentTime)
}
