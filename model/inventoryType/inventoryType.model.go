package model

type GetInventoryTypeResult struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type GetInventoryByIDInput struct {
	ID string `json:"_id" bson:"_id"`
}

type GetInventoryTypeByIDResult struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type AddInventoryTypeInput struct {
	Name string `json:"name" bson:"name"`
}

type UpdateInventoryTypeID struct {
	ID string `json:"_id" bson:"_id"`
}

type UpdateInventoryTypeInput struct {
	Name string `json:"name" bson:"name"`
}

type DeleteInventoryTypeID struct {
	ID string `json:"_id" bson:"_id"`
}
