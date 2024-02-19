package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/location"
	"context"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetLocationService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetLocationResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	pipeline := getPipelineGetLocation(input, searchPipeline)
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetLocationResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetLocationByIDService(input model.GetLocationByIDInput) (model.GetLocationByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetLocationByIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
	var result model.GetLocationByIDResult
	err = ref.FindOne(context.Background(), bson.M{"_id": input.ID, "status": 1}).Decode(&result)
	if err != nil {
		return model.GetLocationByIDResult{}, err
	}
	return result, nil
}

func AddLocationService(input model.AddLocationInput) (string, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return "", err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
	result, err := ref.InsertOne(context.Background(), bson.M{
		"name":    input.Name,
		"address": input.Address,
		"status":  1,
		"created": time.Now(),
	})
	if err != nil {
		return "", err
	}
	return result.InsertedID.(string), nil
}

func UpdateLocationService(input model.UpdateLocationInput, updateID model.UpdateLocationID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	updateField := bson.M{}
	for i := 0; i < reflect.TypeOf(input).NumField(); i++ {
		field := reflect.TypeOf(input).Field(i)
		if !common.IsEmpty(reflect.ValueOf(input).Field(i).Interface()) {
			updateField[field.Tag.Get("bson")] = reflect.ValueOf(input).Field(i).Interface()
		}
	}
	updateField["updatedAt"] = time.Now()
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
	result, err := ref.UpdateOne(context.Background(), bson.M{"_id": updateID.ID, "status": 1}, bson.M{"$set": updateField})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return exception.NotFoundError{Message: "location not found"}
	}
	return nil
}

func DeleteLocationService(deleteID model.DeleteLocationID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
	result, err := ref.UpdateOne(context.Background(), bson.M{"_id": deleteID.ID, "status": 1}, bson.M{"$set": bson.M{"status": 0}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return exception.NotFoundError{Message: "location not found"}
	}
	return nil
}

func GetLocationCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("location")
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

func getPipelineGetLocation(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
	pipeline := bson.A{}
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
		"_id":     1,
		"name":    1,
		"address": 1,
	}}
	pipeline = append(pipeline, matchState, skipState, limitState, sortState, projectState)
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	return pipeline
}
