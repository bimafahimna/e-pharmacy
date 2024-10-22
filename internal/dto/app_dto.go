package dto

type Response struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	TotalRecords int  `json:"total_records,omitempty"`
	TotalPages   int  `json:"total_pages,omitempty"`
	CurrentPage  int  `json:"current_page,omitempty"`
	PreviousPage bool `json:"prev_page"`
	NextPage     bool `json:"next_page"`
}

func PaginationInfo(totalRecords, offset, limit int) *Pagination {
	totalPages := (totalRecords + limit - 1) / limit
	currentPage := offset/limit + 1

	return &Pagination{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
		PreviousPage: currentPage > 1,
		NextPage:     currentPage < totalPages,
	}
}
