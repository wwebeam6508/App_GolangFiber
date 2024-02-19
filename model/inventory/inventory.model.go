package model

import (
	"PBD_backend_go/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetInventoryResult struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	Quantity      int32                `json:"quantity" bson:"quantity"`
	InventoryType entity.InventoryType `json:"inventoryType" bson:"inventoryType"`
}

type GetInventoryByIDInput struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"required"`
}

type GetInventoryByIDResult struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id"`
	Name          string               `json:"name" bson:"name"`
	Description   string               `json:"description" bson:"description"`
	Price         float64              `json:"price" bson:"price"`
	Quantity      int32                `json:"quantity" bson:"quantity"`
	InventoryType entity.InventoryType `json:"inventoryType" bson:"inventoryType"`
}

type AddInventoryInput struct {
	Name          string  `json:"name" bson:"name" validate:"required"`
	Description   string  `json:"description" bson:"description" validate:"required"`
	Price         float64 `json:"price" bson:"price" validate:"required"`
	Quantity      int32   `json:"quantity" bson:"quantity" validate:"required"`
	InventoryType string  `json:"inventoryType" bson:"inventoryType" validate:"required"`
}

type UpdateInventoryID struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"required"`
}

type UpdateInventoryInput struct {
	Name          string  `json:"name" bson:"name"`
	Description   string  `json:"description" bson:"description"`
	Price         float64 `json:"price" bson:"price"`
	Quantity      int32   `json:"quantity" bson:"quantity"`
	InventoryType string  `json:"inventoryType" bson:"inventoryType"`
}

type DeleteInventoryID struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"required"`
}
