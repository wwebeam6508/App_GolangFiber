package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/employee"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func GetEmployeeService(input model.GetEmployeeInput, searchPipeline commonentity.SearchPipeline) ([]model.GetEmployeeResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	pipeline := getPipelineGetEmployee(input, searchPipeline)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetEmployeeResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetEmployeeCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
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
	var result []int32
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0], nil
}

func GetEmployeeByIDService(input model.GetEmployeeByIDInput) (model.GetEmployeeByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetEmployeeByIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	var result model.GetEmployeeByIDResult
	err = ref.FindOne(context.Background(), bson.M{"_id": input.EmployeeID}).Decode(&result)
	if err != nil {
		return model.GetEmployeeByIDResult{}, err
	}
	return result, nil
}

func getPipelineGetEmployee(input model.GetEmployeeInput, searchPipeline commonentity.SearchPipeline) []interface{} {
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
		"employeeID": "$_id",
		"firstName":  1,
		"lastName":   1,
		"bornDate":   0,
		"joinedDate": 1,
		"hiredType":  1,
		"salary":     1}}
	pipeline := bson.A{matchState, skipState, limitState, sortState, projectState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	return pipeline
}
