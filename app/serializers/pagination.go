package serializers

import (
	"ar5go/infra/config"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type ListFilters struct {
	Size int64  `json:"size"`
	Page int64  `json:"page"`
	Sort string `json:"sort"`

	TotalRows    int64  `json:"total_rows"`
	TotalPages   int64  `json:"total_pages"`
	FirstPage    string `json:"first_page"`
	PreviousPage string `json:"previous_page"`
	NextPage     string `json:"next_page"`
	LastPage     string `json:"last_page"`
	FromRow      int64  `json:"from_row"`
	ToRow        int64  `json:"to_row"`

	Results interface{} `json:"results"`

	QueryString string   `json:"qs"`
	Searches    []Search `json:"search"`
}

type Search struct {
	Column string `json:"column,omitempty"`
	Action string `json:"action,omitempty"`
	Query  string `json:"query"`
}

func (lf *ListFilters) GenerateFilters(query url.Values) {
	// default limit, page & sort parameter
	lf.Size = config.App().DefaultPageSize
	lf.Page = 1
	lf.Sort = config.App().Sort
	searchString := ""

	var searches []Search

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "size":
			if size, err := strconv.ParseInt(queryValue, 10, 64); err == nil && size > 0 {
				lf.Size = size
			}
		case "page":
			if page, err := strconv.ParseInt(queryValue, 10, 64); err == nil && page > 0 {
				lf.Page = page
			}
		case "sort":
			lf.Sort = queryValue
		case "qs":
			searchString = queryValue
		}

		// check if query parameter key contains dot
		if strings.Contains(key, ".") {
			// split query parameter key by dot
			searchKeys := strings.Split(key, ".")

			// create search object
			search := Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			// add search object to search array
			searches = append(searches, search)
		}
	}
	lf.QueryString = searchString
	lf.Searches = searches
}

func (lf *ListFilters) GeneratePagesPath(urlPath string) {
	// search query params
	searchQueryParams := ""
	totalPages := lf.TotalPages
	for _, search := range lf.Searches {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	// set first & last page pagination response
	lf.FirstPage = fmt.Sprintf("?size=%d&page=%d&sort=%s", lf.Size, 1, lf.Sort) + searchQueryParams
	lf.LastPage = fmt.Sprintf("?size=%d&page=%d&sort=%s", lf.Size, totalPages, lf.Sort) + searchQueryParams

	if lf.Page > 1 {
		// set previous page pagination response
		lf.PreviousPage = fmt.Sprintf("?size=%d&page=%d&sort=%s", lf.Size, lf.Page-1, lf.Sort) + searchQueryParams
	}

	if lf.Page < totalPages {
		// set next page pagination response
		lf.NextPage = fmt.Sprintf("?size=%d&page=%d&sort=%s", lf.Size, lf.Page+1, lf.Sort) + searchQueryParams
	}

	if lf.Page > totalPages {
		// reset previous page
		lf.PreviousPage = ""
	}

	lf.FirstPage = fmt.Sprintf("%s/%s", urlPath, lf.FirstPage)
	lf.LastPage = fmt.Sprintf("%s/%s", urlPath, lf.LastPage)
	lf.NextPage = fmt.Sprintf("%s/%s", urlPath, lf.NextPage)
	lf.PreviousPage = fmt.Sprintf("%s/%s", urlPath, lf.PreviousPage)
}

func (lf *ListFilters) CalculateTotalPageAndRows(totalRows int64) {
	var totalPages, fromRow, toRow int64 = 0, 0, 0

	// calculate total pages
	totalPages = int64(math.Ceil(float64(totalRows) / float64(lf.Size)))

	if lf.Page == 1 {
		// set from & to row on first page
		fromRow = 1
		toRow = lf.Size
	} else {
		if lf.Page <= totalPages {
			// calculate from & to row
			fromRow = (lf.Page-1)*lf.Size + 1
			toRow = fromRow + lf.Size - 1
		}
	}

	if toRow > totalRows {
		// set to row with total rows
		toRow = totalRows
	}

	lf.FromRow = fromRow
	lf.ToRow = toRow
	lf.TotalPages = totalPages
}
