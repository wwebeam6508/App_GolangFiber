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

type GetEmployeeByIDInput struct {
	EmployeeID string `json:"employeeID" bson:"employeeID" validate:"required"`
}

type GetEmployeeByIDResult struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	JoinedDate time.Time          `json:"joinedDate" bson:"joinedDate"`
	BornDate   time.Time          `json:"bornDate" bson:"bornDate"`
	HiredType  entity.HiredType   `json:"hiredType" bson:"hiredType"`
	Salary     float64            `json:"salary" bson:"salary"`
}

type GetEmployeeResult struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	JoinedDate time.Time          `json:"joinedDate" bson:"joinedDate"`
	HiredType  entity.HiredType   `json:"hiredType" bson:"hiredType"`
	Salary     float64            `json:"salary" bson:"salary"`
}

type AddEmployeeInput struct {
	FirstName  string           `json:"firstName" bson:"firstName" validate:"required"`
	LastName   string           `json:"lastName" bson:"lastName" validate:"required"`
	BornDate   time.Time        `json:"bornDate" bson:"bornDate" validate:"required"`
	JoinedDate time.Time        `json:"joinedDate" bson:"joinedDate" validate:"required"`
	HiredType  entity.HiredType `json:"hiredType" bson:"hiredType" validate:"required"`
	Salary     float64          `json:"salary" bson:"salary" validate:"required"`
}

type UpdateEmployeeID struct {
	EmployeeID string `json:"employeeID" bson:"employeeID" validate:"required"`
}

type UpdateEmployeeInput struct {
	FirstName  string           `json:"firstName" bson:"firstName"`
	LastName   string           `json:"lastName" bson:"lastName"`
	BornDate   time.Time        `json:"bornDate" bson:"bornDate"`
	JoinedDate time.Time        `json:"joinedDate" bson:"joinedDate"`
	HiredType  entity.HiredType `json:"hiredType" bson:"hiredType"`
	Salary     float64          `json:"salary" bson:"salary"`
}

type DeleteEmployeeInput struct {
	EmployeeID string `json:"employeeID" bson:"employeeID" validate:"required"`
}
