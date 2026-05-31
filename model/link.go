package model

type Link struct {
	ID          int64  `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Visits      int64  `json:"visits"`
}
