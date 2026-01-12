package dtos

type PaginationMeta struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int `json:"total_page"`
}