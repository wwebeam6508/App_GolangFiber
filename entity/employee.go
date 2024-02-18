package entity

import (
	"time"
)

type HiredType string

const (
	Contract HiredType = "สัญญาจ้าง"
	FullTime HiredType = "ประจำ"
	Gone     HiredType = "ออก"
)

type Sex string

const (
	Male   Sex = "ชาย"
	Female Sex = "หญิง"
)

var SexOptions = []Sex{Male, Female}

var HiredTypeOptions = []HiredType{Contract, FullTime, Gone}

type Employee struct {
	EmployeeID string    `json:"_id" bson:"_id"`
	FirstName  string    `json:"firstName" bson:"firstName"`
	LastName   string    `json:"lastName" bson:"lastName"`
	BornDate   time.Time `json:"bornDate" bson:"bornDate"`
	JoinedDate time.Time `json:"joinedDate" bson:"joinedDate"`
	HiredType  HiredType `json:"hiredType" bson:"hiredType"`
	Salary     float64   `json:"salary" bson:"salary"`
	Address    string    `json:"address" bson:"address"`
	CitizenID  string    `json:"citizenID" bson:"citizenID"`
	Sex        string    `json:"sex" bson:"sex"`
}
