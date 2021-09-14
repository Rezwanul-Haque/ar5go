package controllers

import (
	"clean/app/serializers"
	"clean/infra/config"
	"clean/infra/errors"
	"clean/infra/logger"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GeneratePaginationRequest(c echo.Context) *serializers.Pagination {
	// default limit, page & sort parameter
	limit := config.App().Limit
	page := config.App().Page
	sort := config.App().Sort
	searchString := ""

	var searches []serializers.Search

	query := c.QueryParams()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "limit":
			limit, _ = strconv.ParseInt(queryValue, 10, 64)
		case "page":
			page, _ = strconv.ParseInt(queryValue, 10, 64)
		case "sort":
			sort = queryValue
		case "qs":
			searchString = queryValue
		}

		// check if query parameter key contains dot
		if strings.Contains(key, ".") {
			// split query parameter key by dot
			searchKeys := strings.Split(key, ".")

			// create search object
			search := serializers.Search{Column: searchKeys[0], Action: searchKeys[1], Query: queryValue}

			// add search object to searches array
			searches = append(searches, search)
		}
	}

	return &serializers.Pagination{Limit: limit, Page: page, Sort: sort, QueryString: searchString, Searches: searches}
}

func GeneratePagesPath(c echo.Context, resp *serializers.Pagination, totalPages int64) {
	// search query params
	searchQueryParams := ""

	for _, search := range resp.Searches {
		searchQueryParams += fmt.Sprintf("&%s.%s=%s", search.Column, search.Action, search.Query)
	}

	// set first & last page pagination response
	resp.FirstPage = fmt.Sprintf("?limit=%d&page=%d&sort=%s", resp.Limit, 1, resp.Sort) + searchQueryParams
	resp.LastPage = fmt.Sprintf("?limit=%d&page=%d&sort=%s", resp.Limit, totalPages, resp.Sort) + searchQueryParams

	if resp.Page > 1 {
		// set previous page pagination response
		resp.PreviousPage = fmt.Sprintf("?limit=%d&page=%d&sort=%s", resp.Limit, resp.Page-1, resp.Sort) + searchQueryParams
	}

	if resp.Page < totalPages {
		// set next page pagination response
		resp.NextPage = fmt.Sprintf("?limit=%d&page=%d&sort=%s", resp.Limit, resp.Page+1, resp.Sort) + searchQueryParams
	}

	if resp.Page > totalPages {
		// reset previous page
		resp.PreviousPage = ""
	}

	urlPath := c.Request().URL.Path

	resp.FirstPage = fmt.Sprintf("%s/%s", urlPath, resp.FirstPage)
	resp.LastPage = fmt.Sprintf("%s/%s", urlPath, resp.LastPage)
	resp.NextPage = fmt.Sprintf("%s/%s", urlPath, resp.NextPage)
	resp.PreviousPage = fmt.Sprintf("%s/%s", urlPath, resp.PreviousPage)
}

func GetUserFromContext(c echo.Context) (*serializers.LoggedInUser, error) {
	user, ok := c.Get("user").(*serializers.LoggedInUser)
	logger.Info(fmt.Sprintf("%+v", user))
	if !ok {
		return nil, errors.ErrNoContextUser
	}

	return user, nil
}
