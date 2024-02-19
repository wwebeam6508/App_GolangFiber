package model

type GetLocationResult struct {
	ID      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type GetLocationByIDInput struct {
	ID string `json:"_id" bson:"_id"`
}

type GetLocationByIDResult struct {
	ID      string `json:"_id" bson:"_id"`
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type AddLocationInput struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type UpdateLocationID struct {
	ID string `json:"_id" bson:"_id"`
}

type UpdateLocationInput struct {
	Name    string `json:"name" bson:"name"`
	Address string `json:"address" bson:"address"`
}

type DeleteLocationID struct {
	ID string `json:"_id" bson:"_id"`
}
