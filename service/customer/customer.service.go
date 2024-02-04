package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/customer"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCustomerService(input model.GetCustomerInput, searchPipeline model.SearchPipeline) ([]model.GetCustomerResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("customers")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "customerID", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "address", Value: 1},
		{Key: "taxID", Value: 1},
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

	var result []model.GetCustomerResult
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

func GetCustomerCountService(searchPipeline model.SearchPipeline, resultChan chan<- int32, errChan chan<- error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		errChan <- err
		return
	}
	ref := coll.Database("PBD").Collection("customers")
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	pipeline := bson.A{matchState}
	if searchPipeline.Search != "" {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline = append(pipeline, groupStage)
	var result []bson.M
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		errChan <- err
		return
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		errChan <- err
		return
	}
	if len(result) == 0 {
		resultChan <- 0
		return
	}
	resultChan <- result[0]["count"].(int32)
}

func GetCustomerByIDService(input model.GetCustomerByIDInput) (model.GetCustomerByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}
	customerIDObjectID, err := primitive.ObjectIDFromHex(input.CustomerID)
	if err != nil {
		return model.GetCustomerByIDResult{}, exception.ValidationError{Message: "invalid customerID"}
	}
	ref := coll.Database("PBD").Collection("customers")
	// aggregate
	matchState := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: customerIDObjectID}, {Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "customerID", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "address", Value: 1},
		{Key: "taxID", Value: 1},
		{Key: "emails", Value: 1},
		{Key: "phones", Value: 1},
	}}}
	pipeline := bson.A{matchState, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}
	var result []model.GetCustomerByIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetCustomerByIDResult{}, err
	}

	//check result empty
	if len(result) == 0 {
		return model.GetCustomerByIDResult{}, exception.NotFoundError{Message: "customer not found"}
	}
	return result[0], nil
}

func AddCustomerService(input model.AddCustomerInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return primitive.NilObjectID, err
	}
	ref := coll.Database("PBD").Collection("customers")
	insertResult, err := ref.InsertOne(context.Background(), bson.D{
		{Key: "name", Value: input.Name},
		{Key: "address", Value: input.Address},
		{Key: "taxID", Value: input.TaxID},
		{Key: "emails", Value: input.Emails},
		{Key: "phones", Value: input.Phones},
		{Key: "status", Value: 1},
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

