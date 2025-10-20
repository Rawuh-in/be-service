package db

import (
	"fmt"
	"math"
	"rawuh-service/internal/shared/model"

	"strings"

	"gorm.io/gorm"
)

func Paginate(value interface{}, v *model.PaginationResponse, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	if v.Limit > 0 || v.Page > 0 {
		var totalRows int64
		db.Model(value).Count(&totalRows)

		v.TotalRows = totalRows
		totalPages := int(math.Ceil(float64(totalRows) / float64(v.Limit)))
		v.TotalPages = int32(totalPages)
	}

	return func(db *gorm.DB) *gorm.DB {
		if v.Limit < 1 || v.Page < 1 {
			return db
		}

		offset := (v.Page - 1) * v.Limit
		if v != nil {
			return db.Limit(int(v.Limit)).Offset(int(offset))
		}
		return db
	}
}

func Sort(s *model.Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if s == nil || s.Column == "" || s.Direction == "" {
			return db
		}
		return db.Order(fmt.Sprintf("%s %s", s.Column, s.Direction))
	}
}

func QueryScoop(v string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if v == "" {
			return db
		}

		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 {
			return db
		}

		columns := strings.Split(parts[0], ",")
		value := parts[1]

		var condition string
		var searchValue string

		switch {
		case strings.HasPrefix(value, "%%"):
			condition = "LIKE"
			searchValue = "%" + value[2:] + "%"
		case strings.HasPrefix(value, "%!"):
			condition = "ILIKE"
			searchValue = "%" + value[2:] + "%"
		default:
			condition = "="
			searchValue = value
		}

		// Build condition: (col1 ILIKE ?) OR (col2 ILIKE ?) ...
		orQuery := db
		for i, col := range columns {
			col = strings.TrimSpace(col)
			clause := fmt.Sprintf("%s %s ?", col, condition)
			if i == 0 {
				orQuery = orQuery.Where(clause, searchValue)
			} else {
				orQuery = orQuery.Or(clause, searchValue)
			}
		}

		return orQuery
	}
}
