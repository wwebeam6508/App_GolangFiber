package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(username string, password string) (LoginServiceResult, error) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.PanicLogging(err)
	}

	matchStage := bson.D{{"$match", bson.D{{"username", username}, {"status", bson.D{{"$eq", 1}}}}}}
	addFieldsStage := bson.D{{"$addFields", bson.D{{"userTypeID", bson.D{{"$toObjectId", "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "userType"}, {"localField", "userTypeID"}, {"foreignField", "_id"}, {"as", "userType"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$userType"}}}}
	projectStage := bson.D{{"$project", bson.D{{"_id", 0}, {"userID", "$_id"}, {"username", 1}, {"password", 1}, {"userType", bson.D{{"userTypeID", "$userType._id"}, {"name", "$userType.name"}, {"permission", "$userType.permission"}}}}}}
	pipeline, err := coll.Database(os.Getenv("MONGO_DB_NAME")).Aggregate(context.Background(), mongo.Pipeline{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage})
	if err != nil {
		exception.PanicLogging(err)
	}
	var result UserProfile
	if err = pipeline.All(context.Background(), &result); err != nil {
		exception.PanicLogging(err)
	}
	isValidPass := bcrypt.CompareHashAndPassword([]byte(result.password), []byte(password))
	if isValidPass == nil {
		exception.PanicLogging(exception.UnauthorizedError{Message: "Invalid password"})
	}
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn

	return LoginServiceResult{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		userProfile:  UserProfile(result),
	}, err
}

type UserType struct {
	userTypeID  string
	name        string
	permissions []string
}

type UserProfile struct {
	userID   string
	username string
	password string
	userType UserType
}

type LoginServiceResult struct {
	accessToken  string
	refreshToken string
	userProfile  UserProfile
}
