package model

type GetInventoryTypeResult struct {
	ID   string `json:"inventoryTypeID" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type GetInventoryByIDInput struct {
	InventoryTypeID string `json:"inventoryTypeID" bson:"inventoryTypeID" validate:"required"`
}

type GetInventoryTypeByIDResult struct {
	ID   string `json:"inventoryTypeID" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type AddInventoryTypeInput struct {
	Name string `json:"name" bson:"name"`
}

type UpdateInventoryTypeID struct {
	InventoryTypeID string `json:"inventoryTypeID" bson:"inventoryTypeID" validate:"required"`
}

type UpdateInventoryTypeInput struct {
	Name string `json:"name" bson:"name"`
}

type DeleteInventoryTypeID struct {
	InventoryTypeID string `json:"inventoryTypeID" bson:"inventoryTypeID" validate:"required"`
}
