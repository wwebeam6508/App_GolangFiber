package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	"PBD_backend_go/model"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserRankByUserIDService(input model.GetUserTypeByUserIDInput) (model.GetUserTypeByUserIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.GetUserTypeByUserIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	userIDObjectID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		return model.GetUserTypeByUserIDResult{}, exception.ValidationError{Message: "invalid userID"}
	}
	//aggregate
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
		{Key: "userTypeID", Value: "$userType._id"},
	}}}

	pipeline := bson.A{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetUserTypeByUserIDResult{}, err
	}
	var result []model.GetUserTypeByUserIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetUserTypeByUserIDResult{}, err
	}
	if len(result) <= 0 {
		return model.GetUserTypeByUserIDResult{}, exception.NotFoundError{Message: "user not found"}
	}
	return result[0], nil
}

func GetUserRankByUserTypeIDService(input model.GetUserRankByUserTypeIDInput) (model.GetUserRankByUserTypeIDResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.GetUserRankByUserTypeIDResult{}, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("userType")
	userTypeIDObjectID, err := primitive.ObjectIDFromHex(input.UserTypeID)
	if err != nil {
		return model.GetUserRankByUserTypeIDResult{}, exception.ValidationError{Message: "invalid userTypeID"}
	}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$eq", Value: userTypeIDObjectID}}}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "rank", Value: 1},
	}}}
	pipeline := bson.A{matchStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return model.GetUserRankByUserTypeIDResult{}, err
	}
	var result []model.GetUserRankByUserTypeIDResult
	err = cursor.All(context.Background(), &result)
	if err != nil {
		return model.GetUserRankByUserTypeIDResult{}, err
	}
	if len(result) <= 0 {
		return model.GetUserRankByUserTypeIDResult{}, exception.NotFoundError{Message: "userType not found"}
	}
	return result[0], nil
}
