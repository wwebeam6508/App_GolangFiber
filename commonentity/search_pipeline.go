package commonentity

import "go.mongodb.org/mongo-driver/bson"

type SearchPipeline struct {
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}
