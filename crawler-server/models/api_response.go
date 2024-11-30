package models

// APIResponse는 API의 표준 응답 구조를 정의합니다.
type APIResponse struct {
	Message string      `json:"message"` // 응답 메시지
	Data    interface{} `json:"data"`    // 실제 응답 데이터
	Error   string      `json:"error"`   // 에러 메시지 (성공 시 빈 문자열)
}
