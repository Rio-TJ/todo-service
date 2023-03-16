package models

import (
	"github.com/SoftclubIT/todo-service/pkg/scopes"
	"gorm.io/gorm"
	"net/url"
	"strconv"
)

type PaginationData struct {
	TotalElements   int64 `json:"total"`
	ElementsPerPage int64 `json:"per_page"`
	CurrentPage     int64 `json:"current_page"`
	LastPage        int64 `json:"last_page"`
	FromElement     int64 `json:"from"`
	ToElement       int64 `json:"to"`
}

func NewPaginationData(DB *gorm.DB, model interface{}, params url.Values) (*PaginationData, error) {
	pd := &PaginationData{}

	result := DB.Model(model).Scopes(scopes.Filter(params)).Count(&pd.TotalElements)
	if result.Error != nil {
		return nil, result.Error
	}

	pd.ElementsPerPage, _ = strconv.ParseInt(params.Get("per_page"), 10, 64)
	if pd.ElementsPerPage == 0 {
		pd.ElementsPerPage = 10
	}

	pd.CurrentPage, _ = strconv.ParseInt(params.Get("page"), 10, 64)
	if pd.CurrentPage == 0 {
		pd.CurrentPage = 1
	}

	pd.LastPage = pd.TotalElements / pd.ElementsPerPage
	if pd.TotalElements&pd.ElementsPerPage != 0 {
		pd.LastPage++
	}

	pd.FromElement = 1 + (pd.CurrentPage-1)*pd.ElementsPerPage

	pd.ToElement = pd.FromElement + pd.ElementsPerPage - 1

	if pd.ToElement > pd.TotalElements {
		pd.ToElement = pd.TotalElements
	}

	if pd.FromElement > pd.ToElement {
		pd.FromElement, pd.ToElement = 0, 0
	}

	return pd, nil
}
