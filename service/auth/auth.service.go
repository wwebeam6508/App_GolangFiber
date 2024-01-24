package service

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/auth"
	"context"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(input model.LoginRequest) model.UserResult {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "username", Value: input.Username}, {Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "userTypeID", Value: bson.D{{Key: "$toObjectId", Value: "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "userType"}, {Key: "localField", Value: "userTypeID"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "userType"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userType"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "userID", Value: "$_id"}, {Key: "username", Value: 1}, {Key: "password", Value: 1}, {Key: "userType", Value: bson.D{{Key: "userTypeID", Value: "$userType._id"}, {Key: "name", Value: "$userType.name"}, {Key: "permission", Value: "$userType.permission"}}}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage})
	if err != nil {
		err := exception.NotFoundError{Message: "user not found"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	var result []model.UserProfileResult
	if err = cursor.All(context.Background(), &result); err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}

	password := result[0].Password
	is_error := bcrypt.CompareHashAndPassword([]byte(*password), []byte(input.Password))
	if is_error != nil {
		err := exception.ValidationError{Message: "invalid password"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	result[0].Password = nil
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn
	claims := jwt.MapClaims{
		"data": result[0],
	}
	accessToken, err := generateJWT(claims)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	// generate refresh token payload as result secrekey is from env signOption issuer audience expiresIn
	refreshToken, err := generateRefreshJWT(claims)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())

	return model.UserResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserProfile:  result[0],
	}
}

func UpdateRefreshTokenService(input model.TokenInput) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}

	filter := bson.D{{Key: "userID", Value: input.UserID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "refreshToken", Value: input.Token}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
}

func CheckRefreshTokenService(input model.TokenInput) bool {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"userID", input.UserID}, {"status", bson.D{{"$eq", 1}}}}}},
		bson.D{{"$project", bson.D{{"_id", 0}, {"refreshToken", 1}}}},
	}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		err := exception.NotFoundError{Message: "user not found"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//check if refresh token is valid
	refreshToken, ok := result[0]["refreshToken"].(string)
	if !ok {
		err := exception.ValidationError{Message: "invalid refresh token"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	split := strings.Split(refreshToken, " ")[1]
	if split != input.Token {
		err := exception.ValidationError{Message: "invalid refresh token"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return true
}

func RefreshTokenService(token string) model.RefreshTokenResult {
	// call verifyJWT function
	claims, err := VerifyJWT(token)
	if err != nil {
		err := exception.ValidationError{Message: "invalid token"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	claimsMap := claims.Claims.(jwt.MapClaims)
	ok := CheckRefreshTokenService(model.TokenInput{
		Token:  token,
		UserID: claimsMap["data"].(map[string]interface{})["userID"].(string),
	})
	if !ok {
		err := exception.ValidationError{Message: "invalid token"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn
	accessToken, err := generateJWT(claimsMap)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	return model.RefreshTokenResult{
		AccessToken: accessToken,
		UserID:      claimsMap["data"].(map[string]interface{})["userID"].(string),
	}
}

func RemoveRefreshTokenService(userID string) error {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return err
	}

	filter := bson.D{{"userID", userID}}
	update := bson.D{{"$set", bson.D{{"refreshToken", ""}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return nil
}

func FetchUserDataService(userID string) model.UserResult {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}

	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: userID},
			{Key: "status", Value: bson.D{
				{Key: "$eq", Value: 1},
			}},
		}},
	}
	addFieldsStage := bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "userTypeID", Value: bson.D{
				{Key: "$toObjectId", Value: "$userTypeID.$id"},
			}},
		}},
	}
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "userType"},
			{Key: "localField", Value: "userTypeID"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "userType"},
		}},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$userType"},
		}},
	}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "userID", Value: "$_id"},
			{Key: "username", Value: 1},
			{Key: "password", Value: 1},
			{Key: "userType", Value: bson.D{
				{Key: "userTypeID", Value: "$userType._id"},
				{Key: "name", Value: "$userType.name"},
				{Key: "permission", Value: "$userType.permission"},
			}},
		}},
	}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage})
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	var result []model.UserProfileResult
	if err = cursor.All(context.Background(), &result); err != nil {
		err := exception.NotFoundError{Message: "user not found"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return model.UserResult{
		UserProfile: result[0],
	}
}

func ChangePasswordDataService(input model.ChangePasswordInput) bool {
	// check password and confirm password
	if input.Password != input.ConfirmPassword {
		//exception errorhandler ValidationError
		err := exception.ValidationError{Message: "password and confirm password must be same"}
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//find userID and status == 1
	filter := bson.D{{"_id", input.UserID}, {"status", bson.D{{"$eq", 1}}}}
	update := bson.D{{"$set", bson.D{{"password", string(hashedPassword)}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		exception.ErrorHandler(&fiber.Ctx{}, err)
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return true
}
