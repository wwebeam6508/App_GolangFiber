package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/wage"
	"context"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetWageService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetWageResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	pipeline := getPipelineGetWage(input, searchPipeline)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetWageResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetWageByIDService(input model.GetWageByIDInput) (*model.GetWageByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	var result model.GetWageByIDResult
	WageIDObject, err := primitive.ObjectIDFromHex(input.WageID)
	if err != nil {
		return nil, err
	}
	//findone
	err = ref.FindOne(context.Background(), bson.M{"_id": WageIDObject, "status": 1}).Decode(&result)
	if err != nil {
		return nil, err
	}
	result.EmployeeName, err = getEmployeeNameService()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func AddWageService(input model.AddWageInput) (model.AddWageResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.AddWageResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	//turn employeeID to ObjectID
	var inputField model.AddWageInputMongo
	inputField.Date = input.Date
	for _, v := range input.Employee {
		employeeID, err := primitive.ObjectIDFromHex(v.EmployeeID)
		if err != nil {
			return model.AddWageResult{}, err
		}
		inputField.Employee = append(inputField.Employee, struct {
			EmployeeID primitive.ObjectID `json:"employeeID" bson:"employeeID"`
			Wage       float64            `json:"wage" bson:"wage"`
		}{EmployeeID: employeeID, Wage: v.Wage})
	}
	inputField.Status = 1

	insertResult, err := ref.InsertOne(context.Background(), inputField)
	if err != nil {
		return model.AddWageResult{}, err
	}
	return model.AddWageResult{WageID: insertResult.InsertedID.(primitive.ObjectID)}, nil
}

func UpdateWageService(inputID model.UpdateWageID, input model.UpdateWageInput) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	//reflect of input get numfield
	input.UpdatedAt = time.Now()
	v := reflect.ValueOf(input)
	updateFields := bson.A{}
	for i := 0; i < v.NumField(); i++ {
		if !common.IsEmpty(v.Field(i).Interface()) {
			// if is removeEmployee or AddEmployee
			if v.Type().Field(i).Name == "RemoveEmployee" || v.Type().Field(i).Name == "AddEmployee" {
				continue
			}
			updateFields = append(updateFields, bson.M{"$set": bson.M{v.Type().Field(i).Tag.Get("json"): v.Field(i).Interface()}})
		}
	}
	WageIDObject, err := primitive.ObjectIDFromHex(inputID.WageID)
	if err != nil {
		return err
	}
	if len(updateFields) != 0 {
		res, err := ref.UpdateOne(context.Background(), bson.M{"_id": WageIDObject, "status": 1}, bson.M{"$set": updateFields})
		if err != nil {
			return err
		}
		if res.MatchedCount == 0 {
			return exception.NotFoundError{Message: "wage not found"}
		}

	}
	//remove employee
	if len(input.RemoveEmployee) > 0 {
		//turn removeEmployee to ObjectID
		removeEmployee := bson.A{}
		for _, v := range input.RemoveEmployee {
			employeeID, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return err
			}
			removeEmployee = append(removeEmployee, employeeID)
		}
		_, err = ref.UpdateOne(context.Background(), bson.M{"_id": WageIDObject, "status": 1}, bson.M{"$pull": bson.M{"employee": bson.M{"employeeID": bson.M{"$in": removeEmployee}}}})
		if err != nil {
			return err
		}
	}
	//add employee
	if len(input.AddEmployee) > 0 {
		addEmployee := bson.A{}
		for _, v := range input.AddEmployee {
			employeeID, err := primitive.ObjectIDFromHex(v.EmployeeID)
			if err != nil {
				return err
			}
			addEmployee = append(addEmployee, bson.M{"employeeID": employeeID, "wage": v.Wage})
		}
		_, err = ref.UpdateOne(context.Background(), bson.M{"_id": WageIDObject, "status": 1}, bson.M{"$push": bson.M{"employee": bson.M{"$each": addEmployee}}})
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteWageService(input model.DeleteWageInput) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	WageIDObject, err := primitive.ObjectIDFromHex(input.WageID)
	if err != nil {
		return err
	}
	_, err = ref.UpdateOne(context.Background(), bson.M{"_id": WageIDObject, "status": 1}, bson.M{"$set": bson.M{"status": 0}})
	if err != nil {
		return err
	}
	return nil
}

func GetWageCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("wage")
	matchStage := bson.M{"$match": bson.M{"status": 1}}
	lookupStage := bson.M{"$lookup": bson.M{"from": "employee", "localField": "employee.employeeID", "foreignField": "_id", "as": "employee"}}
	unwindStage := bson.M{"$unwind": bson.M{"path": "$employee", "preserveNullAndEmptyArrays": true}}
	pipeline := bson.A{matchStage, lookupStage, unwindStage}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline, searchPipeline.SearchPipeline)
	}
	groupStage := bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}
	pipeline = append(pipeline, groupStage)
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
	count := result[0]["count"].(int32)
	return count, nil
}

func getEmployeeNameService() ([]model.GetEmployeeNameResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("employee")
	cursor, err := ref.Find(context.Background(), bson.M{"status": 1})
	if err != nil {
		return nil, err
	}
	var result []model.GetEmployeeNameResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getPipelineGetWage(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
	matchStage := bson.M{"$match": bson.M{"status": 1}}
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
	pipeline := bson.A{matchStage, sortState, skipState, limitState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline, searchPipeline.SearchPipeline)
	}
	projectStage := bson.M{"$project": bson.M{
		"wageID": "$_id",
		//get length of employee array
		"employeeCount": bson.M{"$size": "$employee"},
		"allWages": bson.M{
			"$sum": "$employee.wage",
		},
		"date": 1,
	}}
	pipeline = append(pipeline, projectStage)
	return pipeline
}
