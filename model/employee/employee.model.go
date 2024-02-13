package model

import (
	"PBD_backend_go/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetEmployeeInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type GetEmployeeResult struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	JoinedDate time.Time          `json:"joinedDate" bson:"joinedDate"`
	HiredType  entity.HiredType   `json:"hiredType" bson:"hiredType"`
	Salary     float64            `json:"salary" bson:"salary"`
}
