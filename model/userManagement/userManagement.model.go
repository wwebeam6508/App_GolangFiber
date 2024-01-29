package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserControllerInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type GetUserServiceResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	UserType string             `json:"userType" bson:"userType"`
	Username string             `json:"username" bson:"username"`
	Date     time.Time          `json:"date" bson:"date"`
}

type GetUserServiceInput struct {
	Page           int    `json:"page" bson:"page"`
	PageSize       int    `json:"pageSize" bson:"pageSize"`
	SortTitle      string `json:"sortTitle" bson:"sortTitle"`
	SortType       string `json:"sortType" bson:"sortType"`
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}

type GetUserByIDServiceResult struct {
	UserID   primitive.ObjectID `json:"userID" bson:"userID"`
	Username string             `json:"username" bson:"username"`
	UserType primitive.ObjectID `json:"userType" bson:"userType"`
}

type GetUserByIDInput struct {
	UserID string `json:"userID" bson:"userID"`
}

type AddUserInput struct {
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`
	UserTypeID string `json:"userTypeID" bson:"userTypeID"`
}
