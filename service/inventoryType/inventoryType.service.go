package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/inventoryType"
	"context"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetInventoryTypeService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetInventoryTypeResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	pipeline := getPipelineGetInventory(input, searchPipeline)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetInventoryTypeResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetInventoryTypeByID(input model.GetInventoryByIDInput) (*model.GetInventoryTypeByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	var result []model.GetInventoryTypeByIDResult
	//aggrate
	objectID, _ := primitive.ObjectIDFromHex(input.InventoryTypeID)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}, {Key: "status", Value: 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "inventoryTypeID", Value: "$_id"}, {Key: "name", Value: 1}, {Key: "status", Value: 1}}}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}

func AddInventoryTypeService(input model.AddInventoryTypeInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.ObjectID{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	res, err := ref.InsertOne(context.Background(), bson.D{{Key: "name", Value: input.Name}, {Key: "status", Value: 1}})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if res.InsertedID == nil {
		return primitive.ObjectID{}, exception.ErrorHandler(nil, exception.ValidationError{Message: "Insert Failed"})
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func UpdateInventoryTypeService(input model.UpdateInventoryTypeInput, ID model.UpdateInventoryTypeID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	updateField := bson.A{}
	for i := 0; i < reflect.TypeOf(input).NumField(); i++ {
		field := reflect.TypeOf(input).Field(i)
		value := reflect.ValueOf(input).Field(i).Interface()
		if !common.IsEmpty(value) {
			updateField = append(updateField, bson.M{"$set": bson.M{field.Tag.Get("bson"): value}})
		}
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	objectID, _ := primitive.ObjectIDFromHex(ID.InventoryTypeID)
	res, err := ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateField)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrorHandler(nil, exception.NotFoundError{Message: "Not Found"})
	}
	return nil
}

func DeleteInventoryTypeService(ID model.DeleteInventoryTypeID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	objectID, _ := primitive.ObjectIDFromHex(ID.InventoryTypeID)
	res, err := ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.ErrorHandler(nil, exception.NotFoundError{Message: "Not Found"})
	}
	return nil
}

func GetInventoryTypeCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline := bson.A{matchStage, groupStage}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		//append before group stage
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}
	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0]["count"].(int32), nil
}

func getPipelineGetInventory(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
	matchState := bson.M{"$match": bson.M{"status": 1}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipState := bson.M{"$skip": input.Page * input.PageSize}
	limitState := bson.M{"$limit": input.PageSize}
	var sortValue int
	if !common.IsEmpty(input.SortTitle) && !common.IsEmpty(input.SortType) {
		if input.SortType == "desc" {
			sortValue = -1
		} else {
			sortValue = 1
		}
	}
	sortState := bson.M{"$sort": bson.M{input.SortTitle: sortValue}}
	projectState := bson.M{"$project": bson.M{
		"_id":  "$_id",
		"name": 1,
	}}

	pipeline := bson.A{matchState, skipState, limitState, sortState, projectState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	return pipeline
}
