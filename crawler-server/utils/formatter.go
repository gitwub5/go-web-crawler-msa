package utils

import "github.com/JinHyeokOh01/go-crwl-server/models"

// FormatNoticeMessage는 Notice 데이터를 문자열로 포맷합니다.
func FormatNoticeMessage(notice models.Notice) string {
	return notice.Number + "|" + notice.Title + "|" + notice.Date + "|" + notice.Link
}
