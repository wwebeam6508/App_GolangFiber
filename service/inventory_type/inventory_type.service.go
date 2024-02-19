package inventorytype

import (
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/inventory_type"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func GetInventoryType(input commonentity.PaginateInput) ([]model.GetInventoryTypeResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Database().Client().Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
}

func getPipelineGetEmployee(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) commonentity.SearchPipeline {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}
}
