package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/configuration"
	"PBD_backend_go/entity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/requisition"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRequisitionService(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) ([]model.GetRequisitionResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")
	pipeline := getPipelineGetRequisition(input, searchPipeline)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetRequisitionResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetRequisitionCountService(searchPipeline commonentity.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")
	var result []bson.M
	matchStage := bson.M{"$match": bson.M{"status": 1}}
	groupStage := bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}}
	pipeline := bson.A{matchStage, groupStage}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		//append before groupStage
		pipeline = append(pipeline[:1], append(searchPipeline.SearchPipeline, pipeline[1:]...)...)
	}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return 0, err
	}
	if err = cursor.All(context.Background(), &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}

	return int32(result[0]["count"].(int32)), nil
}

func GetRequisitionByIDService(input model.GetRequisitionByIDInput) (*model.GetRequisitionByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")
	var result []model.GetRequisitionByIDResult
	objectID, _ := primitive.ObjectIDFromHex(input.RequisitionID)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}, {Key: "status", Value: 1}}}}
	projectState := bson.M{"$project": bson.M{
		"requisitionID": "$_id",
		"employeeID":    1,
		"inventries":    1,
		"quantity":      "$inventries.quantity",
		"date":          1,
	}}
	pipeline := bson.A{matchStage, projectState}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, exception.ValidationError{Message: "Requisition not found"}
	}
	result[0].EmployeeOptions, err = GetEmployeeOptionsService()
	if err != nil {
		return nil, err
	}
	result[0].InventoryOptions, err = GetInventoryOptionsService()
	if err != nil {
		return nil, err
	}
	result[0].EndStatusOptions = entity.EndStatusOptions
	return &result[0], nil
}

func AddRequisitionService(input model.AddRequisitionInput) (*model.AddRequisitionResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")
	inventoryRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	//get current quantity of inventory
	for _, v := range input.Inventries {
		// match stage
		objectID, _ := primitive.ObjectIDFromHex(v.InventoryID.Hex())
		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}, {Key: "status", Value: 1}}}}
		// project stage
		projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "quantity", Value: 1}}}}
		pipeline := bson.A{matchStage, projectStage}
		cursor, err := inventoryRef.Aggregate(context.Background(), pipeline)
		if err != nil {
			return nil, err
		}
		var result []bson.M
		if err = cursor.All(context.Background(), &result); err != nil {
			return nil, err
		}
		if len(result) == 0 {
			return nil, exception.ValidationError{Message: "Inventory not found"}
		}
		if result[0]["quantity"].(int32) < v.Quantity {
			return nil, exception.ValidationError{Message: "Inventory not enough"}
		}
	}
	input.CreatedAt = time.Time{}
	input.Status = 1
	result, err := ref.InsertOne(context.Background(), input)
	if err != nil {
		return nil, err
	}
	for _, v := range input.Inventries {
		objectID, _ := primitive.ObjectIDFromHex(v.InventoryID.Hex())
		updateStage := bson.D{{Key: "$inc", Value: bson.D{{Key: "quantity", Value: -v.Quantity}}}}
		_, err = inventoryRef.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateStage)
		if err != nil {
			return nil, err
		}
	}

	return &model.AddRequisitionResult{RequisitionID: result.InsertedID.(primitive.ObjectID)}, nil
}

func UpdateRequisitionStatusService(input model.UpdateRequisitionStatusID, updateInput model.UpdateRequisitionStatusInput) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")

	objectID, _ := primitive.ObjectIDFromHex(input.RequisitionID)
	updateStage := bson.D{{Key: "$set", Value: bson.D{{Key: "endStatus", Value: updateInput.EndStatus}}}}
	res, err := ref.UpdateOne(context.Background(), bson.M{
		"_id":       objectID,
		"endStatus": bson.M{"$eq": ""}},
		updateStage)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return exception.ValidationError{Message: "Update failed"}
	}
	if updateInput.EndStatus == "คืน" {
		inventoryRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
		//update inventory
		//get inventories
		var inventories []bson.M
		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}, {Key: "status", Value: 1}}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "inventries", Value: 1}}}}
		pipeline := bson.A{matchStage, projectStage}
		cursor, err := ref.Aggregate(context.Background(), pipeline)
		if err != nil {
			return err
		}
		if err = cursor.All(context.Background(), &inventories); err != nil {
			return err
		}
		// found or not
		if len(inventories) == 0 {
			return exception.ValidationError{Message: "Requisition not found"}
		}
		//update inventory
		for _, v := range inventories {
			for _, inv := range v["inventries"].(primitive.A) {
				invMap := inv.(primitive.M)
				objectID, _ := primitive.ObjectIDFromHex(invMap["inventoryID"].(primitive.ObjectID).Hex())
				updateStage := bson.D{{Key: "$inc", Value: bson.D{{Key: "quantity", Value: invMap["quantity"].(int32)}}}}
				_, err = inventoryRef.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateStage)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func DeleteRequisitionService(id model.DeleteRequisitionID) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("requisition")
	objectID, _ := primitive.ObjectIDFromHex(id.RequisitionID)
	//update inventory
	//get inventories
	var inventories []bson.M
	matchStage := bson.M{
		"$match": bson.M{
			"_id":    objectID,
			"status": 1,
		},
	}
	projectStage := bson.M{"$project": bson.M{"inventries": 1, "endStatus": 1}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	if err = cursor.All(context.Background(), &inventories); err != nil {
		return err
	}
	if len(inventories) == 0 {
		return exception.ValidationError{Message: "Requisition not found"}
	}
	if !common.IsEmpty(inventories[0]["endStatus"].(string)) {
		return exception.ValidationError{Message: "ใบเบิกนี้ได้ทำการสำเร็จแล้ว ไม่สามารถลบได้"}
	}
	//update inventory
	inventoryRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	for _, v := range inventories {
		for _, inv := range v["inventries"].(primitive.A) {
			invMap := inv.(primitive.M)
			objectID, _ := primitive.ObjectIDFromHex(invMap["inventoryID"].(primitive.ObjectID).Hex())
			updateStage := bson.D{{Key: "$inc", Value: bson.D{{Key: "quantity", Value: invMap["quantity"].(int32)}}}}
			_, err = inventoryRef.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: objectID}}, updateStage)
			if err != nil {
				return err
			}
		}
	}
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

func GetEmployeeOptionsService() ([]model.GetEmployeeNameOptions, error) {
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
	var result []model.GetEmployeeNameOptions
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetInventoryOptionsService() ([]model.GetInventoryNameOptions, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("inventory")
	cursor, err := ref.Find(context.Background(), bson.M{"status": 1})
	if err != nil {
		return nil, err
	}
	var result []model.GetInventoryNameOptions
	if err = cursor.All(context.Background(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getPipelineGetRequisition(input commonentity.PaginateInput, searchPipeline commonentity.SearchPipeline) bson.A {
	matchStage := bson.M{"$match": bson.M{"status": 1}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	lookUpEmployeeIDState := bson.M{"$lookup": bson.M{
		"from":         "employee",
		"localField":   "employeeID",
		"foreignField": "_id",
		"as":           "employeeDetail",
	}}
	unwindEmployeeIDState := bson.M{"$unwind": "$employeeDetail"}
	lookUpInventoryIDState := bson.M{"$lookup": bson.M{
		"from":         "inventory",
		"localField":   "inventries.inventoryID",
		"foreignField": "_id",
		"as":           "inventoryDetail",
	}}
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
	pipeline := bson.A{matchStage, lookUpEmployeeIDState, unwindEmployeeIDState, lookUpInventoryIDState, skipState, limitState, sortState}
	if !common.IsEmpty(searchPipeline.SearchPipeline) {
		pipeline = append(pipeline, searchPipeline.SearchPipeline)
	}
	projectState := bson.M{"$project": bson.M{
		"requisitionID": "$_id",
		"employeeName": bson.M{
			"$concat": bson.A{"$employeeDetail.firstName", " ", "$employeeDetail.lastName"},
		},
		"inventryCount": bson.M{"$size": "$inventries"},
		"date":          1,
	}}
	pipeline = append(pipeline, projectState)
	return pipeline
}
