package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/PuerkitoBio/goquery"
)

const cseURL = "https://ce.khu.ac.kr/ce/user/bbs/BMSR00040/list.do?menuNo=1600045"

// crwlCSENotices는 주어진 URL에서 공지사항 데이터를 크롤링합니다.
func CrwlCSENotices(url string) ([]models.Notice, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var notices []models.Notice

	// HTML에서 데이터 가져오기
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		// '대학' 및 '공지' 글 번호는 제외
		number := strings.TrimSpace(s.Find("td.align-middle").Text())
		if number == "대학" || number == "공지" {
			return
		}
		notice := models.Notice{}
		// 글 번호
		notice.Number = number
		// 제목
		titleLink := s.Find("td.tal a")
		notice.Title = strings.Join(strings.Fields(titleLink.Text()), " ")
		// 날짜
		notice.Date = strings.TrimSpace(s.Find("td:nth-child(4)").Text())
		// 링크 가져오기
		notice.Link = cseURL

		notices = append(notices, notice)
	})

	return notices, nil
}
