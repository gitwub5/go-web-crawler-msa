package crwl

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JinHyeokOh01/go-crwl-server/models"
	"github.com/PuerkitoBio/goquery"
)

// getIDFromURL는 URL에서 글의 고유 ID를 추출합니다.
func getIDFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}
	values, _ := url.ParseQuery(parsedURL.RawQuery)
	return values.Get("wr_id")
}

// crwlSWNotices는 주어진 URL에서 소중단 공지사항 데이터를 크롤링합니다.
func CrwlSWNotices(url string) ([]models.Notice, error) {
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
		notice := models.Notice{}
		// 제목
		titleLink := s.Find(".bo_tit a")
		notice.Title = strings.TrimSpace(titleLink.Text())

		// 링크 가져오기
		if link, exists := titleLink.Attr("href"); exists {
			notice.Link = link
			// 고유 ID 추출
			notice.Number = getIDFromURL(link)
		}

		// 날짜
		notice.Date = strings.TrimSpace(s.Find("td.td_datetime").Text())

		notices = append(notices, notice)
	})

	return notices, nil
}
