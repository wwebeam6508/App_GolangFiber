package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/project"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProjectService(input model.GetProjectInput, searchPipeline model.SearchPipeline) ([]model.GetProjectServiceResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("works")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customer.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customer"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$customer"},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "projectID", Value: "$_id"},
		{Key: "title", Value: 1},
		{Key: "date", Value: bson.D{
			{Key: "$toDate", Value: "$date"},
		}},
		{Key: "profit", Value: 1},
		{Key: "dateEnd", Value: bson.D{
			{Key: "$toDate", Value: "$dateEnd"},
		}},
		{Key: "customer", Value: bson.D{
			{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{
					{Key: "$isArray", Value: "$customer"},
				}},
				{Key: "then", Value: "$customer.name"},
				{Key: "else", Value: ""},
			}}}},
	}}}
	pipeline := bson.A{matchState, lookupStage, unwindStage, projectStage, skipStage, limitStage}
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

	var result []model.GetProjectServiceResult
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

func GetProjectCountService(searchPipeline model.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return 0, err
	}
	ref := coll.Database("PBD").Collection("works")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	pipeline := bson.A{matchState}
	if !common.IsEmpty(searchPipeline.Search) {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline = append(pipeline, groupStage)
	var result []bson.M
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return 0, err
	}
	return result[0]["count"].(int32), nil
}

func GetProjectByIDService(input model.GetProjectByIDInput) (model.GetProjectByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.GetProjectByIDResult{}, err
	}
	ref := coll.Database("PBD").Collection("works")
	pipeline := getPipelineGetProjectByID(input.ProjectID)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetProjectByIDResult{}, err
	}
	var result []model.GetProjectByIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetProjectByIDResult{}, err
	}
	if len(result) == 0 {
		return model.GetProjectByIDResult{}, exception.NotFoundError{Message: "Project not found"}
	}

	return result[0], nil
}

func getPipelineGetProjectByID(projectID string) bson.A {
	projectObjectID, _ := primitive.ObjectIDFromHex(projectID)
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: projectObjectID}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customer.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customer"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$customer"},
	}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "projectID", Value: "$_id"},
		{Key: "title", Value: 1},
		{Key: "date", Value: bson.D{
			{Key: "$toDate", Value: "$date"},
		}},
		{Key: "profit", Value: 1},
		{Key: "dateEnd", Value: bson.D{
			{Key: "$toDate", Value: "$dateEnd"},
		}},
		{Key: "detail", Value: 1},
		{Key: "customer", Value: "$customer._id"},
		{Key: "images", Value: 1}}}}
	return bson.A{matchState, lookupStage, unwindStage, projectStage}
}
