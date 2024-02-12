package commonentity

type PaginationResponse struct {
	Code        int           `json:"code"`
	Message     string        `json:"message"`
	Data        interface{}   `json:"data"`
	CurrentPage int           `json:"currentPage"`
	LastPage    int           `json:"lastPage"`
	Pages       []interface{} `json:"pages"`
}
