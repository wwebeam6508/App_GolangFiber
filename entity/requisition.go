package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EndStatus string

const (
	Return    EndStatus = "คืน"
	Consumed  EndStatus = "ใช้แล้วใช้เลย"
	NotReturn EndStatus = "ไม่คืน"
)

var EndStatusOptions = []EndStatus{Return, Consumed, NotReturn}

type Requisition struct {
	ID         string `json:"_id" bson:"_id"`
	Inventries []struct {
		InventoryID primitive.ObjectID `json:"inventoryID" bson:"inventoryID"`
		Quantity    int                `json:"quantity" bson:"quantity"`
	} `json:"inventroy" bson:"inventroy"`
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	EndStatus  string             `json:"endStatus" bson:"endStatus"`
	Date       time.Time          `json:"date" bson:"date"`
}
