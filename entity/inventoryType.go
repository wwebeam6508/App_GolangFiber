package entity

type InventoryType struct {
	ID   string `json:"_id" bson:"_id"` // Consider using a type like a UUID for IDs
	Name string `json:"name" bson:"name"`
}
