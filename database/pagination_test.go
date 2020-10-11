package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPaginator(t *testing.T) {
	//Test Default Values
	opts := PageOpts{}
	paginator := NewPaginator(opts)

	assert.Equal(t, int64(1), paginator.CurrentPage)
	assert.Equal(t, int64(20), paginator.PerPage)

	//Test Provided Values
	opts.Page = 2
	opts.PerPage = 10
	paginator = NewPaginator(opts)

	assert.Equal(t, int64(2), paginator.CurrentPage)
	assert.Equal(t, int64(10), paginator.PerPage)
}

func TestPaginator_SetOffset(t *testing.T) {
	//Default Case
	opts := PageOpts{}
	paginator := NewPaginator(opts)
	paginator.SetOffset()
	assert.Equal(t, int64(0), paginator.Offset)

	//When a page is provided
	paginator.CurrentPage = 2
	paginator.SetOffset()
	assert.Equal(t, int64(20), paginator.Offset)
}

func TestPaginator_SetTotalPages(t *testing.T) {
	//Default Case
	opts := PageOpts{}
	paginator := NewPaginator(opts)
	paginator.SetTotalPages()
	assert.Equal(t, int64(0), paginator.TotalPages)

	//When a Total Rows is provided
	paginator.TotalRows = 100
	paginator.SetTotalPages()
	assert.Equal(t, int64(5), paginator.TotalPages)

	paginator.TotalRows = 99
	paginator.SetTotalPages()
	assert.Equal(t, int64(5), paginator.TotalPages)
}

func TestPaginator_SetPrevPage(t *testing.T) {
	//Default Case
	opts := PageOpts{}
	paginator := NewPaginator(opts)
	paginator.TotalRows = 100
	paginator.SetPrevPage()
	assert.Equal(t, int64(0), paginator.PrevPage)

	//Test With Current Page > 1
	paginator.CurrentPage = 3
	paginator.SetPrevPage()
	assert.Equal(t, int64(2), paginator.PrevPage)
}

func TestPaginator_SetNextPage(t *testing.T) {
	//Default Case
	opts := PageOpts{}
	paginator := NewPaginator(opts)
	paginator.TotalRows = 100
	paginator.SetNextPage()
	assert.Equal(t, int64(2), paginator.NextPage)

	//Test With Current Page > 1
	paginator.CurrentPage = 5
	paginator.SetNextPage()
	assert.Equal(t, int64(5), paginator.NextPage)
}
