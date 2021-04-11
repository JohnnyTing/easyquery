package easyquery

type Pagination struct {
	Size    int   `json:"size"`
	Page    int   `json:"page"`
	Offset  int   `json:"offset"`
	Pagable bool  `json:"pagable"`
	Total   int64 `json:"total"`
}

func NewPagination(size, page int, pagination bool) *Pagination {
	newPage := &Pagination{
		Size:    size,
		Page:    page,
		Pagable: pagination,
	}
	newPage.Offset = (newPage.Page - 1) * newPage.Size
	return newPage
}

func (p *Pagination) GetSize() int {
	return p.Size
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetOffset() int {
	return p.Offset
}

func (p *Pagination) GetPagable() bool {
	return p.Pagable
}

func (p *Pagination) GetTotal() int64 {
	return p.Total
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}
