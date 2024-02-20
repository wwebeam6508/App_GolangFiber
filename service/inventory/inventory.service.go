package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/inventory"
	"context"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetInventoryService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetInventoryResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	pipeline := getPipelineGetInventory(input, searchPipeline)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetInventoryResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetInventoryByID(input model.GetInventoryByIDInput) (*model.GetInventoryByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	var result []model.GetInventoryByIDResult
	objectID, _ := primitive.ObjectIDFromHex(input.InventoryID)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}, {Key: "status", Value: 1}}}}
	lookUpState := bson.M{"$lookup": bson.M{
		"from":         "inventory_type",
		"localField":   "inventoryType",
		"foreignField": "_id",
		"as":           "inventoryTypeDetail",
	}}
	unwindState := bson.M{"$unwind": "$inventoryTypeDetail"}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "inventoryTypeID", Value: "$_id"}, {Key: "name", Value: 1}, {Key: "description", Value: 1}, {Key: "price", Value: 1}, {Key: "quantity", Value: 1}, {Key: "inventoryType", Value: "$inventoryTypeDetail._id"}}}}
	pipeline := bson.A{matchStage, lookUpState, unwindState, projectStage}
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
	result[0].InventoryTypeName, err = getInventoryTypeNameService()
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func AddInventoryService(input model.AddInventoryInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.NilObjectID, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	inventoryTypeObjectID, _ := primitive.ObjectIDFromHex(input.InventoryType)
	res, err := ref.InsertOne(context.Background(), bson.M{
		"name":          input.Name,
		"description":   input.Description,
		"price":         input.Price,
		"quantity":      input.Quantity,
		"inventoryType": inventoryTypeObjectID,
		"status":        1,
		"createdAt":     time.Now(),
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if res.InsertedID == nil {
		return primitive.NilObjectID, exception.ValidationError{Message: "Insert failed"}
	}
	return res.InsertedID.(primitive.ObjectID), err
}

func UpdateInventoryService(input model.UpdateInventoryInput, id model.UpdateInventoryID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	objectID, _ := primitive.ObjectIDFromHex(id.InventoryID)
	updateField := bson.A{}
	for i := 0; i < reflect.ValueOf(input).NumField(); i++ {
		value := reflect.ValueOf(input).Field(i).Interface()
		field := reflect.TypeOf(input).Field(i)
		if !common.IsEmpty(value) {
			if reflect.TypeOf(input).Field(i).Tag.Get("json") == "inventoryType" {
				inventoryTypeObjectID, _ := primitive.ObjectIDFromHex(input.InventoryType)
				updateField = append(updateField, bson.M{"inventoryType": inventoryTypeObjectID})
				continue
			}
			updateField = append(updateField, bson.M{"$set": bson.M{field.Tag.Get("bson"): value}})
		}
	}
	updateField = append(updateField, bson.M{"$set": bson.M{"updatedAt": time.Now()}})
	res, err := ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateField)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return exception.ValidationError{Message: "Update failed"}
	}
	return nil
}

func DeleteInventoryService(id model.DeleteInventoryID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	objectID, _ := primitive.ObjectIDFromHex(id.InventoryID)
	updateState := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}}
	res, err := ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateState)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return exception.ValidationError{Message: "Delete failed"}
	}
	return nil
}

func GetInventoryCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	matchState := bson.M{"$match": bson.M{"status": 1}}
	groupState := bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}
	pipeline := bson.A{matchState, groupState}
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
	return int32(result[0]["count"].(int32)), nil
}

func getInventoryTypeNameService() ([]model.GetInventoryTypeName, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory_type")
	cursor, err := ref.Find(context.Background(), bson.M{"status": 1})
	if err != nil {
		return nil, err
	}
	var result []model.GetInventoryTypeName
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getPipelineGetInventory(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
	matchState := bson.M{"$match": bson.M{"status": 1}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	lookUpState := bson.M{"$lookup": bson.M{
		"from":         "inventory_type",
		"localField":   "inventoryType",
		"foreignField": "_id",
		"as":           "inventoryTypeDetail",
	}}
	unwindState := bson.M{"$unwind": "$inventoryTypeDetail"}
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
		"_id":           "$_id",
		"name":          1,
		"description":   1,
		"price":         1,
		"quantity":      1,
		"inventoryType": "$inventoryTypeDetail.name",
	}}

	pipeline := bson.A{matchState, lookUpState, unwindState, skipState, limitState, sortState, projectState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline[:2], append(searchPipeline.SearchPipeline, pipeline[2:]...)...)
	}
	return pipeline
}
