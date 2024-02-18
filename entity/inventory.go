package entity

type Inventroy struct {
	ID          string `json:"_id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Price       int    `json:"price" bson:"price"`
	Quantity    int    `json:"quantity" bson:"quantity"`
}
