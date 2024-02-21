package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetInventoryResult struct {
	ID            primitive.ObjectID `json:"inventoryID" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	Quantity      int32              `json:"quantity" bson:"quantity"`
	InventoryType string             `json:"inventoryType" bson:"inventoryType"`
}

type GetInventoryByIDInput struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"required"`
}

type GetInventoryByIDResult struct {
	ID                primitive.ObjectID     `json:"inventoryID" bson:"_id"`
	Name              string                 `json:"name" bson:"name"`
	Description       string                 `json:"description" bson:"description"`
	Price             float64                `json:"price" bson:"price"`
	Quantity          int32                  `json:"quantity" bson:"quantity"`
	InventoryType     primitive.ObjectID     `json:"inventoryType" bson:"inventoryType"`
	InventoryTypeName []GetInventoryTypeName `json:"inventoryTypeName" bson:"inventoryTypeName"`
}

type AddInventoryInput struct {
	Name          string  `json:"name" bson:"name" validate:"required"`
	Description   string  `json:"description" bson:"description" validate:"required"`
	Price         float64 `json:"price,string" bson:"price" validate:"required"`
	Quantity      int32   `json:"quantity,string" bson:"quantity" validate:"required gte=0"`
	InventoryType string  `json:"inventoryType" bson:"inventoryType" validate:"required"`
}

type UpdateInventoryID struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"gte=0"`
}

type UpdateInventoryInput struct {
	Name          string  `json:"name" bson:"name"`
	Description   string  `json:"description" bson:"description"`
	Price         float64 `json:"price,string" bson:"price"`
	Quantity      int32   `json:"quantity,string" bson:"quantity" validate:"gte=0"`
	InventoryType string  `json:"inventoryType" bson:"inventoryType"`
}

type DeleteInventoryID struct {
	InventoryID string `json:"inventoryID" bson:"inventoryID" validate:"required"`
}

type GetInventoryTypeName struct {
	InventoryTypeID primitive.ObjectID `json:"inventoryTypeID" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
}
