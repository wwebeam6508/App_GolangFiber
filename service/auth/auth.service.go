package service

import (
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/auth"
	"context"
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(input model.LoginRequest) (model.UserResult, error) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.UserResult{}, err
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "username", Value: input.Username}, {Key: "status", Value: bson.D{{Key: "$eq", Value: 1}}}}}}
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "userTypeID", Value: bson.D{{Key: "$toObjectId", Value: "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "userType"}, {Key: "localField", Value: "userTypeID"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "userType"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userType"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "userID", Value: "$_id"}, {Key: "username", Value: 1}, {Key: "password", Value: 1}, {Key: "userType", Value: bson.D{{Key: "userTypeID", Value: "$userType._id"}, {Key: "name", Value: "$userType.name"}, {Key: "permission", Value: "$userType.permission"}}}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), mongo.Pipeline{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage})
	if err != nil {
		return model.UserResult{}, err
	}
	var result []model.UserProfileResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return model.UserResult{}, err
	}

	password := result[0].Password
	is_error := bcrypt.CompareHashAndPassword([]byte(*password), []byte(input.Password))
	if is_error != nil {
		return model.UserResult{}, is_error
	}
	result[0].Password = nil
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn
	claims := jwt.MapClaims{
		"data": result[0],
	}
	accessToken, err := generateJWT(claims)
	if err != nil {
		return model.UserResult{}, err
	}
	// generate refresh token payload as result secrekey is from env signOption issuer audience expiresIn
	refreshToken, err := generateRefreshJWT(claims)
	if err != nil {
		return model.UserResult{}, err
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())

	return model.UserResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserProfile:  result[0],
	}, nil
}

func UpdateRefreshTokenService(input model.TokenRequest) error {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return err
	}

	filter := bson.D{{"userID", input.UserID}}
	update := bson.D{{"$set", bson.D{{"refreshToken", input.Token}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return nil
}

func CheckRefreshTokenService(input model.TokenRequest) (bool, error) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return false, err
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"userID", input.UserID}, {"status", bson.D{{"$eq", 1}}}}}},
		bson.D{{"$project", bson.D{{"_id", 0}, {"refreshToken", 1}}}},
	}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return false, err
	}
	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		return false, err
	}
	//check if refresh token is valid
	refreshToken, ok := result[0]["refreshToken"].(string)
	if !ok {
		return false, errors.New("invalid refresh token")
	}
	split := strings.Split(refreshToken, " ")[1]
	if split != input.Token {
		return false, errors.New("invalid refresh token")
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return true, nil
}

func RefreshTokenService(token string) (model.RefreshTokenResult, error) {
	// call verifyJWT function
	claims, err := verifyJWT(token)
	if err != nil {
		return model.RefreshTokenResult{}, err
	}
	claimsMap := claims.Claims.(jwt.MapClaims)
	ok, err := CheckRefreshTokenService(model.TokenRequest{
		Token:  token,
		UserID: claimsMap["data"].(map[string]interface{})["userID"].(string),
	})
	if !ok {
		return model.RefreshTokenResult{}, err
	}
	// generate access token payload as result secrekey is from env signOption issuer audience expiresIn
	accessToken, err := generateJWT(claimsMap)
	if err != nil {
		return model.RefreshTokenResult{}, err
	}
	return model.RefreshTokenResult{
		AccessToken: accessToken,
		UserID:      claimsMap["data"].(map[string]interface{})["userID"].(string),
	}, nil
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

func FetchUserDataService(userID string) (model.UserResult, error) {
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return model.UserResult{}, err
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
		return model.UserResult{}, err
	}
	var result []model.UserProfileResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return model.UserResult{}, err
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return model.UserResult{
		UserProfile: result[0],
	}, nil
}

func ChangePasswordDataService(input model.ChangePasswordRequest) (bool, error) {
	// check password and confirm password
	if input.Password != input.ConfirmPassword {
		return false, errors.New("password and confirm password not match")
	}
	// call the ConnectToMongoDB function
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return false, err
	}
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	//find userID and status == 1
	filter := bson.D{{"_id", input.UserID}, {"status", bson.D{{"$eq", 1}}}}
	update := bson.D{{"$set", bson.D{{"password", string(hashedPassword)}}}}
	collection := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return false, err
	}
	//disconnect from db
	defer coll.Disconnect(context.Background())
	return true, nil
}
