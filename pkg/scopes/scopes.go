package scopes

import (
	"gorm.io/gorm"
	"net/url"
	"strconv"
)

func Paginate(q url.Values) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		perPage, _ := strconv.Atoi(q.Get("per_page"))
		switch {
		case perPage > 100:
			perPage = 100
		case perPage <= 0:
			perPage = 10
		}

		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

func Filter(q url.Values) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		status := q.Get("status")
		if status == "done" || status == "undone" {
			db.Where("status", status)
		}
		return db
	}
}
