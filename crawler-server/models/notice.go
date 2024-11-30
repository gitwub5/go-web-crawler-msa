package models

type Notice struct {
	Number string `json:"number"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Link   string `json:"link"`
}
