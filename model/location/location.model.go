package model

type GetLocationResult struct {
	ID      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type GetLocationByIDInput struct {
	LocationID string `json:"locationID" bson:"locationID" validate:"required"`
}

type GetLocationByIDResult struct {
	ID      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type AddLocationInput struct {
	Name    string `json:"name" bson:"name" validate:"required"`
	Address string `json:"address" bson:"address" validate:"required"`
}

type UpdateLocationID struct {
	LocationID string `json:"locationID" bson:"locationID" validate:"required"`
}

type UpdateLocationInput struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type DeleteLocationID struct {
	LocationID string `json:"locationID" bson:"locationID" validate:"required"`
}
