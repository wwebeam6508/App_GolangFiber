package commonentity

type PaginationResponse struct {
	CurrentPage int           `json:"currentPage"`
	Data        interface{}   `json:"data"`
	LastPage    int           `json:"lastPage"`
	Pages       []interface{} `json:"pages"`
}
