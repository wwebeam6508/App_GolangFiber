package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetWageResult struct {
	WageID        primitive.ObjectID `json:"wageID" bson:"wageID"`
	EmployeeCount int                `json:"employeeCount" bson:"employeeCount"`
	AllWages      int                `json:"allWages" bson:"allWages"`
	Date          time.Time          `json:"date" bson:"date"`
}

type GetWageByIDInput struct {
	WageID string `json:"wageID" bson:"wageID" validate:"required"`
}

type GetWageByIDResult struct {
	WageID   primitive.ObjectID `json:"wageID" bson:"wageID"`
	Employee []struct {
		EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
		Wage       float64            `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee"`
	Date time.Time `json:"date" bson:"date"`
}

type AddWageInput struct {
	Employee []struct {
		EmployeeID string  `json:"employeeID" bson:"employeeID"`
		Wage       float64 `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee" validate:"required"`
	Date string `json:"date" bson:"date" validate:"required"`
}
