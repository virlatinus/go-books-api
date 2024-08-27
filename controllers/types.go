package controllers

type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PaginationResponse struct {
	TotalRecords int  `json:"total_records"`
	CurrentPage  int  `json:"current_page"`
	TotalPages   int  `json:"total_pages"`
	NextPage     *int `json:"next_page"`
	PrevPage     *int `json:"prev_page"`
	PageSize     int  `json:"page_size"`
}

type Response struct {
	Status     StatusResponse      `json:"status"`
	Data       interface{}         `json:"data"`
	Errors     []string            `json:"errors"`
	Pagination *PaginationResponse `json:"pagination"`
}
