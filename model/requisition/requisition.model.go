package model

import (
	"PBD_backend_go/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetRequisitionResult struct {
	RequisitionID primitive.ObjectID `json:"requisitionID" bson:"_id"`
	EmployeeName  string             `json:"employeeName" bson:"employeeName"`
	InventryCount int32              `json:"inventryCount" bson:"inventryCount"`
	EndStatus     string             `json:"endStatus" bson:"endStatus"`
	Date          time.Time          `json:"date" bson:"date"`
}

type GetRequisitionByIDInput struct {
	RequisitionID string `json:"requisitionID" bson:"requisitionID" validate:"required"`
}

type GetRequisitionByIDResult struct {
	RequisitionID primitive.ObjectID `json:"requisitionID" bson:"_id"`
	EmployeeID    primitive.ObjectID `json:"employeeID" bson:"employeeID"`
	Inventries    []struct {
		InventoryID primitive.ObjectID `json:"inventoryID" bson:"inventoryID"`
		Quantity    int32              `json:"quantity" bson:"quantity"`
	} `json:"inventries" bson:"inventries"`
	Date             time.Time                 `json:"date" bson:"date"`
	EndStatus        string                    `json:"endStatus" bson:"endStatus"`
	EmployeeOptions  []GetEmployeeNameOptions  `json:"employeeOptions" bson:"employeeOptions"`
	InventoryOptions []GetInventoryNameOptions `json:"inventoryOptions" bson:"inventoryOptions"`
	EndStatusOptions []entity.EndStatus        `json:"endStatusOptions" bson:"endStatusOptions"`
}

type AddRequisitionInput struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID" validate:"required"`
	Inventries []struct {
		InventoryID primitive.ObjectID `json:"inventoryID" bson:"inventoryID"`
		Quantity    int32              `json:"quantity" bson:"quantity"`
	} `json:"inventries" bson:"inventries" validate:"required"`
	Date      time.Time `json:"date" bson:"date" validate:"required"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Status    int       `json:"status" bson:"status"`
	EndStatus string    `json:"endStatus" bson:"endStatus"`
}

type AddRequisitionResult struct {
	RequisitionID primitive.ObjectID `json:"requisitionID" bson:"requisitionID"`
}

type UpdateRequisitionStatusID struct {
	RequisitionID string `json:"requisitionID" bson:"requisitionID" validate:"required"`
}

type UpdateRequisitionStatusInput struct {
	EndStatus string `json:"endStatus" bson:"endStatus" validate:"required"`
}

type DeleteRequisitionID struct {
	RequisitionID string `json:"requisitionID" bson:"requisitionID" validate:"required"`
}

type GetEmployeeNameOptions struct {
	EmployeeID primitive.ObjectID `json:"employeeID" bson:"_id"`
	FirstName  string             `json:"firstName" bson:"firstName"`
	LastName   string             `json:"lastName" bson:"lastName"`
}

type GetInventoryNameOptions struct {
	InventoryID primitive.ObjectID `json:"inventoryID" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Value       float64            `json:"value" bson:"value"`
}
