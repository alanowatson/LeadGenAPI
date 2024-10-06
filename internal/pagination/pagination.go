package pagination

import (
    "net/http"
    "strconv"
)

const (
    defaultPage    = 1
    defaultPerPage = 10
    maxPerPage     = 100
)

type PaginationParams struct {
    Page    int
    PerPage int
}

func GetPaginationParams(r *http.Request) PaginationParams {
    params := PaginationParams{
        Page:    defaultPage,
        PerPage: defaultPerPage,
    }

    if page := r.URL.Query().Get("page"); page != "" {
        if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
            params.Page = pageNum
        }
    }

    if perPage := r.URL.Query().Get("per_page"); perPage != "" {
        if perPageNum, err := strconv.Atoi(perPage); err == nil && perPageNum > 0 {
            params.PerPage = perPageNum
            if params.PerPage > maxPerPage {
                params.PerPage = maxPerPage
            }
        }
    }

    return params
}

func PaginateSlice(slice interface{}, params PaginationParams) interface{} {
    // Implementation depends on the reflect package to handle different slice types
    // For simplicity, we'll assume we're working with []interface{} for now
    items := slice.([]interface{})

    start := (params.Page - 1) * params.PerPage
    if start > len(items) {
        return []interface{}{}
    }

    end := start + params.PerPage
    if end > len(items) {
        end = len(items)
    }

    return items[start:end]
}
