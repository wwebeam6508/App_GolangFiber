package entity

import (
	"time"
)

type HiredType string

const (
	Contract HiredType = "contract"
	FullTime HiredType = "fulltime"
	Gone     HiredType = "gone"
)

type Employee struct {
	EmployeeID string    `json:"employeeID" bson:"employeeID"`
	FirstName  string    `json:"firstName" bson:"firstName"`
	LastName   string    `json:"lastName" bson:"lastName"`
	BornDate   time.Time `json:"bornDate" bson:"bornDate"`
	JoinedDate time.Time `json:"joinedDate" bson:"joinedDate"`
	HiredType  HiredType `json:"hiredType" bson:"hiredType"`
	Salary     float64   `json:"salary" bson:"salary"`
}
