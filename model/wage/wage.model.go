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
	WageID   primitive.ObjectID `json:"wageID" bson:"_id"`
	Employee []struct {
		EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
		Wage       float64            `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee"`
	Date         time.Time               `json:"date" bson:"date"`
	EmployeeName []GetEmployeeNameResult `json:"employeeName" bson:"employeeName"`
}

type AddWageInput struct {
	Employee []struct {
		EmployeeID string  `json:"employeeID" bson:"employeeID"`
		Wage       float64 `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee" validate:"required"`
	Date time.Time `json:"date" bson:"date" validate:"required"`
}

type AddWageInputMongo struct {
	Employee []struct {
		EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
		Wage       float64            `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee" validate:"required"`
	Date      time.Time `json:"date" bson:"date" validate:"required"`
	Status    int       `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type AddWageResult struct {
	WageID primitive.ObjectID `json:"wageID" bson:"wageID"`
}

type UpdateWageID struct {
	WageID string `json:"wageID" bson:"wageID" validate:"required"`
}

type UpdateWageInput struct {
	AddEmployee []struct {
		EmployeeID string  `json:"employeeID" bson:"employeeID"`
		Wage       float64 `json:"wage" bson:"wage"`
	} `json:"addEmployee" bson:"addEmployee"`
	RemoveEmployee []string  `json:"removeEmployee" bson:"removeEmployee"`
	UpdatedAt      time.Time `json:"updatedAt" bson:"updatedAt"`
}

type DeleteWageInput struct {
	WageID string `json:"wageID" bson:"wageID" validate:"required"`
}

type GetEmployeeNameResult struct {
	EmployeeID string  `json:"employeeID" bson:"_id"`
	FirstName  string  `json:"firstName" bson:"firstName"`
	LastName   string  `json:"lastName" bson:"lastName"`
	Salary     float64 `json:"salary" bson:"salary"`
}
