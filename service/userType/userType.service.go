package service

import (
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/userType"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserTypeService(input model.GetUserTypeInput, searchPipeline model.SearchPipeline) ([]model.GetUserTypeResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("userType")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "userTypeID", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "date", Value: "$createdAt"},
		{Key: "rank", Value: 1},
	}}}
	pipeline := bson.A{matchState, projectStage, skipStage, limitStage}
	if searchPipeline.Search != "" {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	if input.SortTitle != "" && input.SortType != "" {
		var sortValue int
		if input.SortType == "desc" {
			sortValue = -1
		} else {
			sortValue = 1
		}
		sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: input.SortTitle, Value: sortValue}}}}
		pipeline = append(pipeline, sortStage)
	}

	var result []model.GetUserTypeResult
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func GetAllUserTypeCountService(input model.SearchPipeline) int32 {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return 0
	}
	ref := coll.Database("PBD").Collection("userType")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	pipeline := bson.A{matchState}
	if input.Search != "" {
		pipeline = append(pipeline, input.SearchPipeline...)
	}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline = append(pipeline, groupStage)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0
	}
	var result []bson.M
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return 0
	}
	return result[0]["count"].(int32)
}
