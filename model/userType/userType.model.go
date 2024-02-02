package model

import (
	"PBD_backend_go/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserTypeInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type GetUserTypeResult struct {
	UserTypeID primitive.ObjectID `json:"userTypeID" bson:"userTypeID"`
	Name       string             `json:"name" bson:"name"`
	Date       time.Time          `json:"date" bson:"date"`
	Rank       int32              `json:"rank" bson:"rank"`
}

type GetUserTypeByIDInput struct {
	UserTypeID string `json:"userTypeID" bson:"userTypeID"`
}

type GetUserTypeByIDResult struct {
	UserTypeID primitive.ObjectID `json:"userTypeID" bson:"userTypeID"`
	Name       string             `json:"name" bson:"name"`
	Permission entity.Permissions `json:"permission" bson:"permission"`
	Rank       int32              `json:"rank" bson:"rank"`
}

type SearchPipeline struct {
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}

type AddUserTypeInput struct {
	Name       string             `json:"name" bson:"name"`
	Permission entity.Permissions `json:"permission" bson:"permission"`
	Rank       int32              `json:"rank" bson:"rank"`
}

type UpdateUserTypeID struct {
	UserTypeID string `json:"userTypeID" bson:"userTypeID"`
}

type UpdateUserTypeInput struct {
	Name       string             `json:"name" bson:"name"`
	Permission entity.Permissions `json:"permission" bson:"permission"`
	Rank       int32              `json:"rank" bson:"rank"`
}

type DeleteUserTypeID struct {
	UserTypeID string `json:"userTypeID" bson:"userTypeID"`
}
