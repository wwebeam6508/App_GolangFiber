package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/user"
	"context"
	"errors"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GetUserService(input model.GetUserServiceInput, searchPipeline model.SearchPipeline) ([]model.GetUserServiceResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("users")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	if input.Page > 0 {
		input.Page = input.Page - 1
	}
	skipStage := bson.D{{Key: "$skip", Value: input.Page * input.PageSize}}
	limitStage := bson.D{{Key: "$limit", Value: input.PageSize}}
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "userTypeID", Value: bson.D{{Key: "$toObjectId", Value: "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "userType"},
		{Key: "localField", Value: "userTypeID"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "userType"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userType"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "userID", Value: "$_id"},
		{Key: "userType", Value: "$userType.name"},
		{Key: "username", Value: 1},
		{Key: "date", Value: "$createdAt"},
	}}}
	pipeline := bson.A{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage, skipStage, limitStage}
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
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetUserServiceResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetUserByIDService(input model.GetUserByIDInput) (model.GetUserByIDServiceResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.GetUserByIDServiceResult{}, err
	}
	userIDObjectID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return model.GetUserByIDServiceResult{}, exception.ValidationError{Message: "invalid userID"}
	}
	ref := coll.Database("PBD").Collection("users")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userIDObjectID}, {Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "userTypeID", Value: bson.D{{Key: "$toObjectId", Value: "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "userType"},
		{Key: "localField", Value: "userTypeID"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "userType"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userType"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "userID", Value: "$_id"},
		{Key: "username", Value: 1},
		{Key: "userTypeID", Value: "$userType._id"},
	}}}
	pipeline := bson.A{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetUserByIDServiceResult{}, err
	}
	// check is cursor empty
	var result []model.GetUserByIDServiceResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetUserByIDServiceResult{}, err
	}
	// check is result empty
	if len(result) <= 0 {
		return model.GetUserByIDServiceResult{}, exception.NotFoundError{Message: "user not found"}
	}
	return result[0], nil
}

func AddUserService(input model.AddUserInput) (primitive.ObjectID, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return primitive.ObjectID{}, err
	}
	userTypeIDObjectID, err := primitive.ObjectIDFromHex(input.UserTypeID)
	if err != nil {
		return primitive.ObjectID{}, exception.ValidationError{Message: "invalid userTypeID"}
	}
	password, err := encryptPassword(input.Password)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	ref := coll.Database("PBD").Collection("users")
	ires, err := ref.InsertOne(context.Background(), bson.D{
		{Key: "username", Value: input.Username},
		{Key: "password", Value: password},
		{Key: "userTypeID", Value: bson.D{{Key: "$ref", Value: "userType"}, {Key: "$id", Value: userTypeIDObjectID}}},
		{Key: "createdAt", Value: primitive.NewDateTimeFromTime(time.Now())},
		{Key: "status", Value: 1},
	})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if ires.InsertedID == nil {
		return primitive.ObjectID{}, errors.New("failed to add user")
	}

	return ires.InsertedID.(primitive.ObjectID), nil
}

func UpdateUserService(input model.UpdateUserInput, id model.UpdateUserID) error {

	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return err
	}
	//get only input that is not empty dynamic
	updateData := bson.D{}
	reflectInput := reflect.ValueOf(input)
	for i := 0; i < reflectInput.NumField(); i++ {
		// if paasword or selfID or userID or userTypeID then continue
		if reflectInput.Type().Field(i).Name == "Password" || reflectInput.Type().Field(i).Name == "SelfID" || reflectInput.Type().Field(i).Name == "UserTypeID" {
			continue
		}

		if reflectInput.Field(i).Interface() != "" {
			updateData = append(updateData, bson.E{Key: reflectInput.Type().Field(i).Tag.Get("bson"), Value: reflectInput.Field(i).Interface()})
		}
	}
	// encode password
	if input.Password != "" {
		password, err := encryptPassword(input.Password)
		if err != nil {
			return err
		}
		updateData = append(updateData, bson.E{Key: "password", Value: password})
	}
	//userTypeID go for ref
	if input.UserTypeID != "" {
		userTypeIDObjectID, err := primitive.ObjectIDFromHex(input.UserTypeID)
		if err != nil {
			return exception.ValidationError{Message: "invalid userTypeID"}
		}
		updateData = append(updateData, bson.E{Key: "userTypeID", Value: bson.D{{Key: "$ref", Value: "userType"}, {Key: "$id", Value: userTypeIDObjectID}}})
	}

	userIDObjectID, err := primitive.ObjectIDFromHex(id.UserID)
	if err != nil {
		return exception.ValidationError{Message: "invalid userID"}
	}
	coll.Database("PBD").Collection("users").UpdateOne(context.Background(), bson.D{{Key: "_id", Value: userIDObjectID}}, bson.D{{Key: "$set", Value: updateData}})
	return nil
}

func DeleteUserService(input model.DeleteUserInput) error {

	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return err
	}
	userIDObjectID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return exception.ValidationError{Message: "invalid userID"}
	}
	coll.Database("PBD").Collection("users").UpdateOne(context.Background(), bson.D{{Key: "_id", Value: userIDObjectID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: 0}}}})
	return nil
}

func GetAllUserCount(searchPipeline model.SearchPipeline) int32 {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return 0
	}
	ref := coll.Database("PBD").Collection("users")
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}
	pipeline := bson.A{matchStage, groupStage}
	if searchPipeline.Search != "" {
		pipeline = append(pipeline, searchPipeline.SearchPipeline...)
	}
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

func GetUserTypeService() ([]model.GetUserTypeServiceResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return nil, err
	}
	ref := coll.Database("PBD").Collection("userType")
	//aggregate
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: "$_id"},
		{Key: "name", Value: 1},
		{Key: "rank", Value: 1},
	}}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	var result []model.GetUserTypeServiceResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to encrypt password")
	}
	return string(hash), nil
}
