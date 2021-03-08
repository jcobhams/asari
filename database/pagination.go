package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

var (
	DefaultPageNumber  int64 = 1
	DefaultPerPageRows int64 = 20
)

type (
	Paginator struct {
		CurrentPage int64 `json:"currentPage"`
		NextPage    int64 `json:"nextPage"`
		PrevPage    int64 `json:"prevPage"`
		TotalPages  int64 `json:"totalPages"`
		TotalRows   int64 `json:"totalRows"`
		PerPage     int64 `json:"perPage"`
		Offset      int64 `json:"-"`
	}

	PageOpts struct {
		Page    int64 `json:"page"`
		PerPage int64 `json:"per_page"`
	}

	PaginatedResult struct {
		Paginator
		Cursor *mongo.Cursor
	}
)

func NewPaginator(opts PageOpts) *Paginator {
	p := &Paginator{}
	p.Offset = 0

	if opts.Page <= 1 {
		p.CurrentPage = DefaultPageNumber
	} else {
		p.CurrentPage = opts.Page
	}

	if opts.PerPage < 1 {
		p.PerPage = DefaultPerPageRows
	} else {
		p.PerPage = opts.PerPage
	}

	return p
}

func (p *Paginator) SetOffset() {
	if p.CurrentPage == 1 {
		p.Offset = 0
		return
	}
	p.Offset = (p.CurrentPage - 1) * p.PerPage
}

func (p *Paginator) SetTotalPages() {
	if p.TotalRows == 0 {
		p.TotalPages = 0
		return
	}
	p.TotalPages = int64(math.Ceil(float64(p.TotalRows) / float64(p.PerPage)))
}

func (p *Paginator) SetPrevPage() {
	//Call SetTotalPages just to be safe
	p.SetTotalPages()

	if p.CurrentPage == 1 {
		p.PrevPage = 0
		return
	}
	p.PrevPage = p.CurrentPage - 1
}

func (p *Paginator) SetNextPage() {
	//Call SetTotalPages just to be safe
	p.SetTotalPages()

	if p.CurrentPage == p.TotalPages {
		p.NextPage = p.CurrentPage
	} else {
		p.NextPage = p.CurrentPage + 1
	}
}
