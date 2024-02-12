package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/expense"
	"context"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetExpenseService(input model.GetExpenseInput, searchPipeline model.SearchPipeline) (result []model.GetExpenseServiceResult, err error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")

	result = []model.GetExpenseServiceResult{}
	cursor, err := ref.Aggregate(context.Background(), getPipelineGetExpense(input, searchPipeline))
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetExpenseByIDService(input model.GetExpenseByIDInput) (model.GetExpenseByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetExpenseByIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	pipeline := getPipelineGetExpenseByID(input)
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetExpenseByIDResult{}, err
	}
	var result []model.GetExpenseByIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetExpenseByIDResult{}, err
	}
	if len(result) == 0 {
		return model.GetExpenseByIDResult{}, exception.NotFoundError{Message: "Expense Not found"}
	}
	return result[0], nil
}

func AddExpenseService(input model.AddExpenseInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.NilObjectID, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	// add status
	input.Status = 1

	for i := 0; i < len(input.Lists); i++ {
		input.Lists[i].ID = primitive.NewObjectID()
	}
	addInput := bson.D{}
	for i := 0; i < reflect.ValueOf(input).NumField(); i++ {
		if !common.IsEmpty(reflect.ValueOf(input).Field(i).Interface()) {
			if reflect.ValueOf(input).Type().Field(i).Name == "WorkRef" || reflect.ValueOf(input).Type().Field(i).Name == "CustomerRef" {
				objectID, _ := primitive.ObjectIDFromHex(reflect.ValueOf(input).Field(i).Interface().(string))
				addInput = append(addInput, bson.E{Key: reflect.ValueOf(input).Type().Field(i).Tag.Get("json"), Value: bson.D{{Key: "$ref", Value: reflect.ValueOf(input).Type().Field(i).Tag.Get("json")}, {Key: "$id", Value: objectID}}})
				continue
			}
			addInput = append(addInput, bson.E{Key: reflect.ValueOf(input).Type().Field(i).Tag.Get("json"), Value: reflect.ValueOf(input).Field(i).Interface()})
		}
	}

	res, err := ref.InsertOne(context.Background(), addInput)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func UpdateExpenseService(input model.UpdateExpenseInput, expenseID string) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	expenseIDObjectID, _ := primitive.ObjectIDFromHex(expenseID)

	//check input is empty or not
	updateInput := bson.D{}
	for i := 0; i < reflect.ValueOf(input).NumField(); i++ {
		if !common.IsEmpty(reflect.ValueOf(input).Field(i).Interface()) {
			// exclude addLists and removeLists
			if reflect.ValueOf(input).Type().Field(i).Name == "AddLists" || reflect.ValueOf(input).Type().Field(i).Name == "RemoveLists" {
				continue
			}
			if reflect.ValueOf(input).Type().Field(i).Name == "WorkRef" || reflect.ValueOf(input).Type().Field(i).Name == "CustomerRef" {
				objectID, _ := primitive.ObjectIDFromHex(reflect.ValueOf(input).Field(i).Interface().(string))
				updateInput = append(updateInput, bson.E{Key: reflect.ValueOf(input).Type().Field(i).Tag.Get("json"), Value: bson.D{{Key: "$ref", Value: reflect.ValueOf(input).Type().Field(i).Tag.Get("json")}, {Key: "$id", Value: objectID}}})
				continue
			}
			updateInput = append(updateInput, bson.E{Key: reflect.ValueOf(input).Type().Field(i).Tag.Get("json"), Value: reflect.ValueOf(input).Field(i).Interface()})
			continue
		}
	}

	filter := bson.D{{Key: "_id", Value: expenseIDObjectID}}
	updateFirst := bson.D{{Key: "$set", Value: updateInput}}
	_, err = ref.UpdateOne(context.Background(), filter, updateFirst)
	if err != nil {
		return err
	}

	//addLists and removeLists
	if len(input.AddLists) > 0 {
		//add ID to addLists
		for i := 0; i < len(input.AddLists); i++ {
			input.AddLists[i].ID = primitive.NewObjectID()
		}
		//add addLists to lists
		updateAddLists := bson.D{{Key: "$push", Value: bson.D{{Key: "lists", Value: bson.D{{Key: "$each", Value: input.AddLists}}}}}}
		_, err = ref.UpdateOne(context.Background(), filter, updateAddLists)

		if err != nil {
			return err
		}
	}
	if len(input.RemoveLists) > 0 {
		//remove lists get only their ID
		removeListID := []primitive.ObjectID{}
		for i := 0; i < len(input.RemoveLists); i++ {
			//
			RemoveListsID, _ := primitive.ObjectIDFromHex(input.RemoveLists[i].ID)
			removeListID = append(removeListID, RemoveListsID)
		}
		//remove removeLists from lists
		updateRemoveLists := bson.D{{Key: "$pull", Value: bson.D{{Key: "lists", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: removeListID}}}}}}}}
		_, err = ref.UpdateOne(context.Background(), filter, updateRemoveLists)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleleExpenseService(input model.DeleteExpenseInput) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	expenseIDObjectID, _ := primitive.ObjectIDFromHex(input.ExpenseID)
	filter := bson.D{{Key: "_id", Value: expenseIDObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}}
	_, err = ref.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetWorkTitleService() ([]model.GetWorkTitle, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "title", Value: 1},
		{Key: "id", Value: "$_id"},
	}}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetWorkTitle
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCustomerNameService() ([]model.GetCustomerName, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "name", Value: 1},
		{Key: "id", Value: "$_id"},
	}}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetCustomerName
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getPipelineGetExpense(input model.GetExpenseInput, searchPipeline model.SearchPipeline) bson.A {
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
	//unwind stage allow empty array
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$workRef"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	lookupStageCustomer := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customerRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customerRef"},
	}}}
	//unwind stage allow empty array
	unwindStageCustomer := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$customerRef"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		//expenseID
		{Key: "expenseID", Value: "$_id"},
		{Key: "title", Value: 1},
		{Key: "date", Value: bson.D{
			{Key: "$toDate", Value: "$date"},
		}},
		{Key: "isVat", Value: bson.D{
			{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$currentVat", 0}}}, false, true,
			}},
		}},
		{Key: "totalPrice", Value: bson.D{
			{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$lists"},
				{Key: "initialValue", Value: 0},
				{Key: "in", Value: bson.D{
					{Key: "$add", Value: bson.A{"$$value", "$$this.price"}},
				}}}},
		}},
		{Key: "workRef", Value: "$workRef.title"},
		{Key: "customerRef", Value: "$customerRef.name"},
	}}}

	pipeline := bson.A{matchStage, lookupStage, lookupStageCustomer, unwindStage, unwindStageCustomer, projectStage, skipStage, limitStage}
	if !common.IsEmpty(searchPipeline.Search) && len(searchPipeline.SearchPipeline) > 0 {
		//put searchPipeline to pipeline before projectStage
		pipeline = append(pipeline[:3], append(searchPipeline.SearchPipeline, pipeline[3:]...)...)
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
	return pipeline
}

func getPipelineGetExpenseByID(input model.GetExpenseByIDInput) bson.A {
	expenseIDObjectID, _ := primitive.ObjectIDFromHex(input.ExpenseID)
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: expenseIDObjectID}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "works"},
		{Key: "localField", Value: "workRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "workRef"},
	}}}
	// unwind stage allow empty array
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$workRef"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	lookupStageCustomer := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customerRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customerRef"},
	}}}
	// unwind stage allow empty array
	unwindStageCustomer := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$customerRef"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "expenseID", Value: "$_id"},
		{Key: "title", Value: 1},
		{Key: "date", Value: bson.D{
			{Key: "$toDate", Value: "$date"},
		}},
		{Key: "lists", Value: 1},
		{Key: "currentVat", Value: 1},
		{Key: "detail", Value: 1},
		{Key: "workRef", Value: "$workRef._id"},
		{Key: "customerRef", Value: "$customerRef._id"},
	}}}
	pipeline := bson.A{matchStage, lookupStage, lookupStageCustomer, unwindStage, unwindStageCustomer, projectStage}
	return pipeline
}

func GetExpenseCountService(searchPipeline model.SearchPipeline) (int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "works"},
		{Key: "localField", Value: "workRef.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "workRef"}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline := bson.A{matchStage, lookupStage, groupStage}
	if !common.IsEmpty(searchPipeline.Search) {
		pipeline = append(pipeline[:2], append(searchPipeline.SearchPipeline, pipeline[2:]...)...)
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
		return 0, nil
	}
	return (result[0]["count"].(int32)), nil
}
