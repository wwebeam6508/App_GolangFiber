package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/project"
	"context"
	"os"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProjectService(input model.GetProjectInput, searchPipeline model.SearchPipeline) ([]model.GetProjectServiceResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
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
		{Key: "preserveNullAndEmptyArrays", Value: true},
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
		{Key: "customer", Value: "$customer.name"},
	}}}
	pipeline := bson.A{matchState, lookupStage, unwindStage, projectStage, skipStage, limitStage}
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
	defer coll.Disconnect(context.Background())
	if err != nil {
		return 0, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "customers"},
		{Key: "localField", Value: "customer.$id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "customer"},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline := bson.A{matchStage, lookupStage, groupStage}
	if !common.IsEmpty(searchPipeline.Search) && len(searchPipeline.SearchPipeline) > 0 {
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
	return result[0]["count"].(int32), nil
}

func GetProjectByIDService(input model.GetProjectByIDInput) (model.GetProjectByIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return model.GetProjectByIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
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
	return result[0], nil
}

func AddProjectService(input model.AddProjectInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return primitive.NilObjectID, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	//exclude images
	addInput := bson.D{}
	for i := 0; i < reflect.ValueOf(input).NumField(); i++ {
		if reflect.ValueOf(input).Type().Field(i).Name != "Images" {
			if reflect.ValueOf(input).Type().Field(i).Name == "Customer" {
				objectID, _ := primitive.ObjectIDFromHex(input.Customer)
				addInput = append(addInput, bson.E{Key: "customer", Value: bson.D{{Key: "$id", Value: objectID}}})
				continue
			}
			addInput = append(addInput, bson.E{Key: reflect.ValueOf(input).Type().Field(i).Tag.Get("json"), Value: reflect.ValueOf(input).Field(i).Interface()})
		}
	}
	insertResult, err := ref.InsertOne(context.Background(), addInput)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func UpdateProjectService(input model.UpdateProjectInput, projectID string) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	projectObjectID, _ := primitive.ObjectIDFromHex(projectID)

	update := bson.D{}
	//reflect input that is not empty
	inputRef := reflect.ValueOf(input)
	for i := 0; i < inputRef.NumField(); i++ {
		if !common.IsEmpty(inputRef.Field(i).Interface()) {
			if inputRef.Type().Field(i).Name == "Images" {
				continue
			}
			if inputRef.Type().Field(i).Name == "Customer" {
				objectID, _ := primitive.ObjectIDFromHex(inputRef.Field(i).Interface().(string))
				update = append(update, bson.E{Key: "customer", Value: bson.D{{Key: "$id", Value: objectID}}})
				continue
			}
			update = append(update, bson.E{Key: inputRef.Type().Field(i).Tag.Get("json"), Value: inputRef.Field(i).Interface()})
		}
	}
	_, err = ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: projectObjectID}}, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return err
	}
	return nil
}

func DeleteProjectService(projectID string) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	projectObjectID, _ := primitive.ObjectIDFromHex(projectID)
	_, err = ref.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: projectObjectID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}})
	if err != nil {
		return err
	}
	return nil
}

func GetCustomerNameService() ([]model.GetCustomerNameResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("customers")
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: "$_id"},
		{Key: "name", Value: 1},
	}}}
	pipeline := bson.A{projectStage}
	var result []model.GetCustomerNameResult
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

func ForceDeleteProjectService(projectID string) error {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	projectObjectID, _ := primitive.ObjectIDFromHex(projectID)
	_, err = ref.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: projectObjectID}})
	if err != nil {
		return err
	}
	return nil
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
		{Key: "preserveNullAndEmptyArrays", Value: true},
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
