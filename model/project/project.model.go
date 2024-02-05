package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProjectInput struct {
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

type GetProjectServiceResult struct {
	ProjectID primitive.ObjectID `json:"projectID" bson:"projectID"`
	Title     string             `json:"title" bson:"title"`
	Date      time.Time          `json:"date" bson:"date"`
	Profit    float64            `json:"profit" bson:"profit"`
	DateEnd   time.Time          `json:"dateEnd" bson:"dateEnd"`
	Customer  string             `json:"customer" bson:"customer"`
}

type GetProjectByIDInput struct {
	ProjectID string `json:"projectID" bson:"projectID" validate:"required"`
}

type GetProjectByIDResult struct {
	ProjectID primitive.ObjectID `json:"projectID" bson:"projectID"`
	Title     string             `json:"title" bson:"title"`
	Date      time.Time          `json:"date" bson:"date"`
	Profit    float64            `json:"profit" bson:"profit"`
	DateEnd   time.Time          `json:"dateEnd" bson:"dateEnd"`
	Detail    string             `json:"detail" bson:"detail"`
	Customer  primitive.ObjectID `json:"customer" bson:"customer"`
	Images    []string           `json:"images" bson:"images"`
}
