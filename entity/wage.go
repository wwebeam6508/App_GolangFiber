package entity

import "time"

type Wage struct {
	WageID   string `json:"wageID" bson:"wageID"`
	Employee []struct {
		EmployeeID string  `json:"employeeID" bson:"employeeID"`
		Wage       float64 `json:"wage" bson:"wage"`
	} `json:"employee" bson:"employee"`
	Date time.Time `json:"date" bson:"date"`
}
