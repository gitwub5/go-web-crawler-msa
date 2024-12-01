package crwl

import (
	"net/http"
	"strings"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/PuerkitoBio/goquery"
)

// CrwlCSENotices는 주어진 URL에서 컴퓨터공학과 공지사항 데이터를 크롤링합니다.
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
		number := strings.TrimSpace(s.Find("td.align-middle").Text())
		if number == "대학" || number == "공지" {
			return
		}
		notice := models.Notice{
			Number: number,
			Title:  strings.Join(strings.Fields(s.Find("td.tal a").Text()), " "),
			Date:   strings.TrimSpace(s.Find("td:nth-child(4)").Text()),
			Link:   url,
		}
		notices = append(notices, notice)
	})

	return notices, nil
}
