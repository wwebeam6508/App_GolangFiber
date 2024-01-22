package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/auth"
	"context"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(input model.LoginRequest) (bson.M, error) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.PanicLogging(err)
	}

	matchStage := bson.D{{"$match", bson.D{{"username", input.Username}, {"status", bson.D{{"$eq", 1}}}}}}
	addFieldsStage := bson.D{{"$addFields", bson.D{{"userTypeID", bson.D{{"$toObjectId", "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "userType"}, {"localField", "userTypeID"}, {"foreignField", "_id"}, {"as", "userType"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$userType"}}}}
	projectStage := bson.D{{"$project", bson.D{{"_id", 0}, {"userID", "$_id"}, {"username", 1}, {"password", 1}, {"userType", bson.D{{"userTypeID", "$userType._id"}, {"name", "$userType.name"}, {"permission", "$userType.permission"}}}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage})
	if err != nil {
		exception.PanicLogging(err)
	}
	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		exception.PanicLogging(err)
	}

	password := result[0]["password"].(string)
	is_error := bcrypt.CompareHashAndPassword([]byte(password), []byte(input.Password))
	if is_error != nil {
		return bson.M{}, is_error
	}
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn
	claims := jwt.MapClaims{
		"userProfile": result[0],
	}
	accessToken, err := generateJWT(claims)
	if err != nil {
		exception.PanicLogging(err)
	}
	// generate refresh token payload as result secrekey is from env signOption issuer audience expiresIn
	refreshToken, err := generateRefreshJWT(claims)
	if err != nil {
		exception.PanicLogging(err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())

	return bson.M{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"userProfile":  result[0],
	}, err
}
