package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/expense"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func GetExpenseService(input model.GetExpenseInput, searchPipeline model.SearchPipeline) (result []model.GetExpenseServiceResult, err error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("expenses")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "works"},
		{Key: "localField", Value: "workRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "workRef"},
	}}}
	lookupStageCustomer := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customerRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customerRef"},
	}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "expenseID", Value: "$_id"},
		{Key: "title", Value: 1},
		{Key: "date", Value: bson.D{
			{Key: "$toDate", Value: "$date"},
		}},
		{Key: "lists", Value: 1},
		{Key: "currentVat", Value: 1},
		//turn customerRef and workRef to string
		{Key: "workRef", Value: bson.D{
			{Key: "$arrayElemAt", Value: bson.A{"$workRef", 0}},
		}},
		{Key: "customerRef", Value: bson.D{
			{Key: "$arrayElemAt", Value: bson.A{"$customerRef", 0}},
		}},
	}}}

	pipeline := bson.A{matchStage, lookupStage, lookupStageCustomer, projectStage, skipStage, limitStage}
	if !common.IsEmpty(searchPipeline.Search) && len(searchPipeline.SearchPipeline) > 0 {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	if !common.IsEmpty(input.SortTitle) && !common.IsEmpty(input.SortType) {
		var sortValue int
		if input.SortType == "desc" {
			sortValue = -1
		} else {
			sortValue = 1
		}
		sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: input.SortTitle, Value: sortValue}}}}
		pipeline = append(pipeline, sortStage)
	}
	fmt.Println(pipeline)
	result = []model.GetExpenseServiceResult{}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	coll.Disconnect(context.Background())
	return result, nil
}

func GetExpenseCountService(searchPipeline model.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database("PBD").Collection("expenses")

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline := bson.A{matchStage, groupStage}
	if !common.IsEmpty(searchPipeline.Search) {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	var result []bson.M
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, exception.NotFoundError{Message: "Not found"}
	}
	return (result[0]["count"].(int32)), nil
}
