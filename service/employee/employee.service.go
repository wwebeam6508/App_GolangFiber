package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/employee"
	"context"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetEmployeeService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetEmployeeResult, error) {
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
	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	return result[0]["count"].(int32), nil
}

func GetEmployeeByIDService(input model.GetEmployeeByIDInput) (model.GetEmployeeByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetEmployeeByIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	employeeIDObjectID, _ := primitive.ObjectIDFromHex(input.EmployeeID)
	var result []model.GetEmployeeByIDResult
	matchStage := bson.M{"$match": bson.M{"_id": employeeIDObjectID, "status": 1}}
	projectStage := bson.M{"$project": bson.M{
		"employeeID": "$_id",
		"firstName":  1,
		"lastName":   1,
		"joinedDate": 1,
		"bornDate":   1,
		"hiredType":  1,
		"salary":     1,
		"address":    1,
		"citizenID":  1,
		"sex":        1,
	}}
	cursor, err := ref.Aggregate(context.Background(), bson.A{matchStage, projectStage})
	if err != nil {
		return model.GetEmployeeByIDResult{}, err
	}

	if err = cursor.All(context.Background(), &result); err != nil {
		return model.GetEmployeeByIDResult{}, err
	}
	if len(result) == 0 {
		return model.GetEmployeeByIDResult{}, exception.NotFoundError{Message: "employee not found"}
	}
	return result[0], nil
}

func AddEmployeeService(input model.AddEmployeeInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.ObjectID{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	input.Status = 1
	input.CreatedAt = time.Now()
	res, err := ref.InsertOne(context.Background(), input)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if res.InsertedID == nil {
		return primitive.ObjectID{}, exception.ValidationError{Message: "insertedID is nil"}
	}
	return res.InsertedID.(primitive.ObjectID), nil

}

func UpdateEmployeeService(input model.UpdateEmployeeInput, employee model.UpdateEmployeeID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	employeeIDObjectID, _ := primitive.ObjectIDFromHex(employee.EmployeeID)
	input.UpdatedAt = time.Now()
	updateStage := bson.A{}
	for i := 0; i < reflect.TypeOf(input).NumField(); i++ {
		field := reflect.TypeOf(input).Field(i)
		value := reflect.ValueOf(input).Field(i).Interface()
		if !common.IsEmpty(value) {
			updateStage = append(updateStage, bson.M{"$set": bson.M{field.Tag.Get("bson"): value}})
		}
	}
	res, err := ref.UpdateOne(context.Background(), bson.M{"_id": employeeIDObjectID, "status": 1}, updateStage)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.NotFoundError{Message: "employee not found"}
	}
	return nil
}

func DeleteEmployeeService(employee model.UpdateEmployeeID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	employeeIDObjectID, _ := primitive.ObjectIDFromHex(employee.EmployeeID)
	updateStage := bson.A{bson.M{"$set": bson.M{"status": 0}}}
	res, err := ref.UpdateOne(context.Background(), bson.M{"_id": employeeIDObjectID, "status": 1}, updateStage)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return exception.NotFoundError{Message: "employee not found"}
	}
	return nil
}

func getPipelineGetEmployee(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
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
		"joinedDate": 1,
		"hiredType":  1,
		"salary":     1,
		"sex":        1,
	}}
	pipeline := bson.A{matchState, skipState, limitState, sortState, projectState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	return pipeline
}
