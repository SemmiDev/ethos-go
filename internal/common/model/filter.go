package model

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Filter represents comprehensive filter options for list queries
type Filter struct {
	// Pagination
	CurrentPage int `json:"current_page" form:"current_page" query:"current_page"` // current page (berpindah-pindah halaman)
	PerPage     int `json:"per_page" form:"per_page" query:"per_page"`             // limit (batas data yang ditampilkan)

	// Search
	Keyword string `json:"keyword" form:"keyword" query:"keyword"` // search keyword (keyword pencarian)

	// Sorting
	SortBy        string `json:"sort_by" form:"sort_by" query:"sort_by"`                      // column name to sort
	SortDirection string `json:"sort_direction" form:"sort_direction" query:"sort_direction"` // asc or desc direction

	// Date Range Filters
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date" query:"start_date"` // filter from date
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date" query:"end_date"`       // filter to date

	// Status Filters
	IsActive   *bool `json:"is_active,omitempty" form:"is_active" query:"is_active"`       // filter by active status
	IsInactive *bool `json:"is_inactive,omitempty" form:"is_inactive" query:"is_inactive"` // filter by inactive status
}

// ListResponse is a generic wrapper for paginated list responses
// This ensures consistent response structure across all list endpoints
func NewFilter() Filter {
	return Filter{
		CurrentPage:   DefaultCurrentPage,
		PerPage:       DefaultPageLimit,
		SortDirection: DefaultColumnDirection,
	}
}

const (
	AscDirection           string = "asc"
	DescDirection          string = "desc"
	DefaultPageLimit       int    = 20
	DefaultCurrentPage     int    = 1
	DefaultColumnDirection string = AscDirection
	UnlimitedPage          int    = -1
)

func (f *Filter) GetLimit() int {
	return f.PerPage
}

func (f *Filter) GetOffset() int {
	// ex:
	// page 1 -> (1 - 1) * 10 = 0
	// page 2 -> (2 - 1) * 10 = 10
	offset := (f.CurrentPage - 1) * f.PerPage
	return offset
}

func (f *Filter) HasKeyword() bool {
	return f.Keyword != ""
}

func (f *Filter) HasSort() bool {
	return f.SortBy != ""
}

func (f *Filter) IsDesc() bool {
	return strings.EqualFold(f.SortDirection, DescDirection)
}

func (f *Filter) IsUnlimitedPage() bool {
	return isUnlimitedPage(f.PerPage)
}

func isUnlimitedPage(perPage int) bool {
	return perPage == UnlimitedPage
}

type Paging struct {
	HasPreviousPage        bool `json:"has_previous_page"`
	HasNextPage            bool `json:"has_next_page"`
	CurrentPage            int  `json:"current_page"`
	PerPage                int  `json:"per_page"`
	TotalData              int  `json:"total_data"`
	TotalDataInCurrentPage int  `json:"total_data_in_current_page"`
	LastPage               int  `json:"last_page"`
	From                   int  `json:"from"`
	To                     int  `json:"to"`
}

var ErrPaging = errors.New("per_page harus lebih besar dari 0 dan offset tidak boleh negatif")

func NewPaging(currentPage, perPage, totalData int) (*Paging, error) {
	if isUnlimitedPage(perPage) {
		return &Paging{
			CurrentPage:            currentPage,
			PerPage:                perPage,
			TotalData:              totalData,
			LastPage:               1,
			From:                   1,
			To:                     totalData,
			TotalDataInCurrentPage: totalData,
		}, nil
	}

	if totalData == 0 {
		return &Paging{
			HasPreviousPage:        false,
			HasNextPage:            false,
			CurrentPage:            1,
			PerPage:                perPage,
			TotalData:              0,
			TotalDataInCurrentPage: 0,
			LastPage:               1,
			From:                   0,
			To:                     0,
		}, nil
	}

	offset := (currentPage - 1) * perPage

	if perPage <= 0 || offset < 0 {
		return nil, ErrPaging
	}

	lastPage := totalData / perPage
	if totalData%perPage != 0 {
		lastPage++
	}

	to := min(offset+perPage, totalData)
	from := int(0)
	if to > offset {
		from = offset + 1
	}

	if currentPage > lastPage {
		currentPage = lastPage
	}

	totalDataInCurrentPage := to - offset

	return &Paging{
		HasPreviousPage:        currentPage > 1,
		HasNextPage:            currentPage < lastPage,
		CurrentPage:            currentPage,
		PerPage:                perPage,
		TotalData:              totalData,
		LastPage:               lastPage,
		From:                   from,
		To:                     to,
		TotalDataInCurrentPage: totalDataInCurrentPage,
	}, nil
}

// FilterFromRequest parses filter parameters from an HTTP request
func FilterFromRequest(r *http.Request) Filter {
	query := r.URL.Query()

	filter := NewFilter()

	// Parse pagination
	if page := query.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.CurrentPage = p
		}
	}
	if perPage := query.Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 {
			filter.PerPage = pp
		}
	}

	// Allow "limit" as alias for "per_page"
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.PerPage = l
		}
	}

	// Parse search keyword
	filter.Keyword = strings.TrimSpace(query.Get("keyword"))
	if filter.Keyword == "" {
		filter.Keyword = strings.TrimSpace(query.Get("search"))
	}
	if filter.Keyword == "" {
		filter.Keyword = strings.TrimSpace(query.Get("q"))
	}

	// Parse sorting
	filter.SortBy = strings.TrimSpace(query.Get("sort_by"))
	if filter.SortBy == "" {
		filter.SortBy = strings.TrimSpace(query.Get("order_by"))
	}

	filter.SortDirection = strings.ToLower(strings.TrimSpace(query.Get("sort_direction")))
	if filter.SortDirection == "" {
		filter.SortDirection = strings.ToLower(strings.TrimSpace(query.Get("order")))
	}
	if filter.SortDirection != AscDirection && filter.SortDirection != DescDirection {
		filter.SortDirection = DefaultColumnDirection
	}

	// Parse date range
	if startDate := query.Get("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}
	if endDate := query.Get("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	// Parse status filters
	if active := query.Get("active"); active != "" {
		b := active == "true" || active == "1"
		filter.IsActive = &b
	}
	if inactive := query.Get("inactive"); inactive != "" {
		b := inactive == "true" || inactive == "1"
		filter.IsInactive = &b
	}

	return filter
}

// Validate validates the filter and sets defaults for invalid values
func (f *Filter) Validate() {
	if f.CurrentPage < 1 {
		f.CurrentPage = DefaultCurrentPage
	}
	if f.PerPage < 1 && f.PerPage != UnlimitedPage {
		f.PerPage = DefaultPageLimit
	}
	if f.SortDirection != AscDirection && f.SortDirection != DescDirection {
		f.SortDirection = DefaultColumnDirection
	}
}

// AllowedSortColumns validates if the sort column is in the allowed list
func (f *Filter) ValidateSortBy(allowedColumns []string) bool {
	if f.SortBy == "" {
		return true
	}
	for _, col := range allowedColumns {
		if strings.EqualFold(f.SortBy, col) {
			f.SortBy = col // normalize to expected case
			return true
		}
	}
	return false
}

// ActiveOnly returns true if only active items should be returned
func (f *Filter) ActiveOnly() bool {
	return f.IsActive != nil && *f.IsActive && (f.IsInactive == nil || !*f.IsInactive)
}

// InactiveOnly returns true if only inactive items should be returned
func (f *Filter) InactiveOnly() bool {
	return f.IsInactive != nil && *f.IsInactive && (f.IsActive == nil || !*f.IsActive)
}
