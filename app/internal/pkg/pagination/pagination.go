package pagination

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type Page struct {
	Page   int `json:"page"`
	Size   int `json:"size"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (p *Page) Compute() {
	if p.Page == 0 {
		p.Page = 1
	}

	if p.Size == 0 {
		p.Size = 1
	}

	p.Limit = p.Size
	p.Offset = ((p.Page - 1) * p.Size)

}

func ParsePaginationRequest(r *http.Request) (Page, error) {
	var pagination Page
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	size := r.URL.Query().Get("size")
	if size == "" {
		size = "10"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return pagination, fmt.Errorf("failed to parse page params: %v", err)
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		return pagination, fmt.Errorf("failed to parse size params: %v", err)
	}

	pagination.Page = pageInt
	pagination.Size = sizeInt

	return pagination, nil
}

type Metadata struct {
	CurrentPage  int   `json:"current_page" example:"1"`
	PageSize     int   `json:"page_size" example:"50"`
	FirstPage    int   `json:"first_page" example:"1"`
	LastPage     int   `json:"last_page" example:"1"`
	TotalRecords int64 `json:"total_records" example:"1"`
}

func (m *Metadata) Compute(totalRecords int64, size int, currentPage int) {
	m.TotalRecords = totalRecords
	m.PageSize = size
	m.FirstPage = 1
	m.CurrentPage = currentPage
	m.LastPage = int(math.Ceil(float64(m.TotalRecords) / float64(m.PageSize)))
}
