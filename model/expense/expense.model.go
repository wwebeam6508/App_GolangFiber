package model

import (
	"PBD_backend_go/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type GetExpenseInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type SearchPipeline struct {
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}

type GetExpenseServiceResult struct {
	ExpenseID   string        `json:"expenseID" bson:"expenseID"`
	Title       string        `json:"title" bson:"title"`
	Date        time.Time     `json:"date" bson:"date"`
	Lists       []entity.List `json:"lists" bson:"lists"`
	CurrentVat  int           `json:"currentVat" bson:"currentVat"`
	WorkRef     string        `json:"workRef" bson:"workRef"`
	CustomerRef string        `json:"customerRef" bson:"customerRef"`
}
