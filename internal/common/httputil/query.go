package httputil

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/semmidev/ethos-go/internal/common/model"
)

// GetStringQuery GetString extracts a string query param (trimmed)
func GetStringQuery(r *http.Request, key string, defaultValue string) string {
	v := strings.TrimSpace(r.URL.Query().Get(key))
	if v == "" {
		return defaultValue
	}
	return v
}

// GetIntQuery GetInt extracts an int query param (with default fallback)
func GetIntQuery(r *http.Request, key string, defaultValue int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}
	return i
}

// GetBoolPtrQuery GetBool extracts a bool query param (with default fallback)
func GetBoolPtrQuery(r *http.Request, key string) *bool {
	v := strings.ToLower(strings.TrimSpace(r.URL.Query().Get(key)))
	if v == "" {
		return nil
	}

	b := false
	if v == "true" || v == "1" || v == "yes" {
		b = true
	}
	return &b
}

// ParseFilterQuery ParseFilter parses pagination, sorting, and keyword filters from httputil
func ParseFilterQuery(r *http.Request) model.Filter {
	return model.Filter{
		CurrentPage:   GetIntQuery(r, "current_page", model.DefaultCurrentPage),
		PerPage:       GetIntQuery(r, "per_page", model.DefaultPageLimit),
		Keyword:       GetStringQuery(r, "keyword", ""),
		SortBy:        GetStringQuery(r, "sort_by", "created_at"),
		SortDirection: GetStringQuery(r, "sort_direction", model.DefaultColumnDirection),
	}
}
