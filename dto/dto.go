package dto

import "link-storage-service/model"

type CreateLinkRequest struct {
	URL string `json:"url"`
}

type CreateLinkResponse struct {
	ShortCode string `json:"short_code"`
}

type ListLinksResponse struct {
	Items  []LinkItem `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

type LinkItem struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Visits    int64  `json:"visits"`
}

type GetLinkResponse struct {
	URL    string `json:"url"`
	Visits int64  `json:"visits"`
}

type LinkStatsResponse struct {
	ShortCode string `json:"short_code"`
	URL       string `json:"url"`
	Visits    int64  `json:"visits"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func LinkItemFromModel(link model.Link) LinkItem {
	return LinkItem{
		ID:        link.ID,
		ShortCode: link.ShortCode,
		URL:       link.OriginalURL,
		CreatedAt: link.CreatedAt,
		Visits:    link.Visits,
	}
}

func LinkItemsFromModel(links []model.Link) []LinkItem {
	items := make([]LinkItem, len(links))
	for i, link := range links {
		items[i] = LinkItemFromModel(link)
	}
	return items
}
