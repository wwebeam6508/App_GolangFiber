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
	CurrentVat  float64       `json:"currentVat" bson:"currentVat"`
	WorkRef     string        `json:"workRef" bson:"workRef"`
	CustomerRef string        `json:"customerRef" bson:"customerRef"`
}

type GetExpenseByIDInput struct {
	ExpenseID string `json:"expenseID" bson:"expenseID" validate:"required"`
}

type GetExpenseByIDResult struct {
	ExpenseID   string        `json:"expenseID" bson:"expenseID"`
	Title       string        `json:"title" bson:"title"`
	Date        time.Time     `json:"date" bson:"date"`
	Lists       []entity.List `json:"lists" bson:"lists"`
	CurrentVat  float64       `json:"currentVat" bson:"currentVat"`
	Detail      string        `json:"detail" bson:"detail"`
	WorkRef     string        `json:"workRef" bson:"workRef"`
	CustomerRef string        `json:"customerRef" bson:"customerRef"`
}

type AddExpenseInput struct {
	Title       string        `json:"title" bson:"title" validate:"required"`
	Date        time.Time     `json:"date" bson:"date" validate:"required"`
	Lists       []entity.List `json:"lists" bson:"lists" validate:"required"`
	CurrentVat  float64       `json:"currentVat" bson:"currentVat" validate:"required"`
	Detail      string        `json:"detail" bson:"detail"`
	WorkRef     string        `json:"workRef" bson:"workRef"`
	CustomerRef string        `json:"customerRef" bson:"customerRef" `
	Status      int           `json:"status" bson:"status"`
}

type UpdateExpenseID struct {
	ExpenseID string `json:"expenseID" bson:"expenseID" validate:"required"`
}
type UpdateExpenseInput struct {
	Title       string        `json:"title" bson:"title"`
	Date        time.Time     `json:"date" bson:"date"`
	AddLists    []entity.List `json:"addLists" bson:"addLists"`
	RemoveLists []RemoveList  `json:"removeLists" bson:"removeLists"`
	CurrentVat  float64       `json:"currentVat" bson:"currentVat"`
	Detail      string        `json:"detail" bson:"detail"`
	WorkRef     string        `json:"workRef" bson:"workRef"`
	CustomerRef string        `json:"customerRef" bson:"customerRef"`
}

type RemoveList struct {
	ID string `json:"_id" bson:"_id"`
}
