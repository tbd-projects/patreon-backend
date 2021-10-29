package models

type Pagination struct {
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

