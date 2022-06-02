package db

import (
	"ar5go/app/serializers"
	"ar5go/app/utils/methodsutil"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func TransactionStart() *gorm.DB {
	return Client().DB.Begin()
}

func TransactionComplete(tx interface{}) {
	tx.(*gorm.DB).Commit()
}

func applyFilteringCondition(stmt *gorm.DB, tableName string, filters *serializers.ListFilters, forCount bool) *gorm.DB {
	offset := (filters.Page - 1) * filters.Size
	sort := filters.Sort

	if !methodsutil.IsEmpty(tableName) {
		sort = tableName + "." + filters.Sort
	}

	find := stmt

	// get data with limit, offset & order
	if !forCount {
		find = stmt.Limit(int(filters.Size)).Offset(int(offset)).Order(sort)
	}

	// generate where query
	searches := filters.Searches

	if searches != nil {
		for _, value := range searches {
			var column string
			column = value.Column
			if !methodsutil.IsEmpty(tableName) {
				column = tableName + "." + value.Column
			}
			action := value.Action
			query := value.Query

			switch action {
			case "equals":
				whereQuery := fmt.Sprintf("%s = ?", column)
				find = find.Where(whereQuery, query)
			case "contains":
				whereQuery := fmt.Sprintf("%s LIKE ?", column)
				find = find.Where(whereQuery, "%"+query+"%")
			case "in":
				whereQuery := fmt.Sprintf("%s IN (?)", column)
				queryArray := methodsutil.StringToIntArray(strings.Split(query, ","))
				find = find.Where(whereQuery, queryArray)
			case "gt":
				whereQuery := fmt.Sprintf("%s > (?)", column)
				queryArray := query
				find = find.Where(whereQuery, queryArray)
			case "gte":
				whereQuery := fmt.Sprintf("%s >= (?)", column)
				queryArray := query
				find = find.Where(whereQuery, queryArray)
			case "lt":
				whereQuery := fmt.Sprintf("%s < (?)", column)
				queryArray := query
				find = find.Where(whereQuery, queryArray)
			case "lte":
				whereQuery := fmt.Sprintf("%s <= (?)", column)
				queryArray := query
				find = find.Where(whereQuery, queryArray)
			}
		}
	}

	return find
}

func applyQueryStringSearch(stmt *gorm.DB, searchStmt, qs string) *gorm.DB {
	searchTerm := "%" + qs + "%"
	stmt.Where(searchStmt, map[string]interface{}{"st": searchTerm})

	return stmt
}
