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
	EmployeeID       primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	FirstName        string             `json:"firstName" bson:"firstName"`
	LastName         string             `json:"lastName" bson:"lastName"`
	JoinedDate       time.Time          `json:"joinedDate" bson:"joinedDate"`
	BornDate         time.Time          `json:"bornDate" bson:"bornDate"`
	HiredType        entity.HiredType   `json:"hiredType" bson:"hiredType"`
	Salary           float64            `json:"salary" bson:"salary"`
	Address          string             `json:"address" bson:"address"`
	CitizenID        string             `json:"citizenID" bson:"citizenID"`
	Sex              string             `json:"sex" bson:"sex"`
	HiredTypeOptions []entity.HiredType `json:"hiredTypeOptions" bson:"hiredTypeOptions"`
	SexOptions       []entity.Sex       `json:"sexOptions" bson:"sexOptions"`
}

type GetEmployeeResult struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
	JoinedDate time.Time          `json:"joinedDate" bson:"joinedDate"`
	HiredType  entity.HiredType   `json:"hiredType" bson:"hiredType"`
	Salary     float64            `json:"salary" bson:"salary"`
	Sex        string             `json:"sex" bson:"sex"`
}

type AddEmployeeInput struct {
	FirstName  string           `json:"firstName" bson:"firstName" validate:"required"`
	LastName   string           `json:"lastName" bson:"lastName" validate:"required"`
	BornDate   time.Time        `json:"bornDate" bson:"bornDate" validate:"required"`
	JoinedDate time.Time        `json:"joinedDate" bson:"joinedDate" validate:"required"`
	HiredType  entity.HiredType `json:"hiredType" bson:"hiredType" validate:"required"`
	Sex        string           `json:"sex" bson:"sex" validate:"required"`
	CitizenID  string           `json:"citizenID" bson:"citizenID" validate:"required"`
	Address    string           `json:"address" bson:"address"`
	Salary     float64          `json:"salary" bson:"salary"`
	Status     int              `json:"status" bson:"status"`
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
	Sex        string           `json:"sex" bson:"sex"`
	CitizenID  string           `json:"citizenID" bson:"citizenID"`
	Address    string           `json:"address" bson:"address"`
	Salary     float64          `json:"salary" bson:"salary"`
}

type DeleteEmployeeInput struct {
	EmployeeID string `json:"employeeID" bson:"employeeID" validate:"required"`
}
