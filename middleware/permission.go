package middleware

import (
	"PBD_backend_go/configuration"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model"
	service "PBD_backend_go/service/auth"
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Permission(c *fiber.Ctx, input model.PermissionInput) error {
	//get Authorization from header
	split := strings.Split(c.Get("Authorization"), " ")
	if len(split) != 2 {
		return exception.ErrorHandler(c, exception.UnauthorizedError{Message: "permission denied"})
	}
	token := split[1]
	//call verify jwt
	claim, err := service.VerifyJWT(token)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	//get userID from jwt
	jwtMap := claim.Claims.(jwt.MapClaims)
	userID := jwtMap["data"].(map[string]interface{})["userID"].(string)
	//call check permission
	permission, err := checkPermissionByUserID(userID, input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	if !permission {
		return exception.ErrorHandler(c, exception.UnauthorizedError{Message: "permission denied"})
	} else {
		return c.Next()
	}
}

func checkPermissionByUserID(userID string, input model.PermissionInput) (bool, error) {
	coll, err := configuration.ConnectToMongoDB()
	if err != nil {
		return false, err
	}
	ref := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("users")
	userIDObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userIDObjectID}}}}
	addFieldsStage := bson.D{{Key: "$addFields", Value: bson.D{{Key: "userTypeID", Value: bson.D{{Key: "$toObjectId", Value: "$userTypeID.$id"}}}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "userType"},
		{Key: "localField", Value: "userTypeID"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "userType"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$userType"}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "userType", Value: bson.D{
			{Key: "permission", Value: 1},
		}},
	}}}
	pipeline := bson.A{matchStage, addFieldsStage, lookupStage, unwindStage, projectStage}
	cursor, err := ref.Aggregate(context.Background(), pipeline)
	if err != nil {
		return false, err
	}
	// check is cursor empty
	var result []model.PermissionResult
	if err = cursor.All(context.Background(), &result); err != nil {
		return false, err
	}
	// check is result empty
	if len(result) <= 0 {
		return false, exception.NotFoundError{Message: "user not found"}
	}
	// check permission input has group and name
	if input.GroupName == "" || input.Name == "" {
		return false, exception.ValidationError{Message: "invalid permission input"}
	}
	permission := result[0].UserType.Permission
	val := reflect.ValueOf(permission)
	for i := 0; i < val.NumField(); i++ {
		if val.Type().Field(i).Name == input.GroupName {
			permissionDetail := reflect.ValueOf(val.Field(i).Interface())
			for j := 0; j < permissionDetail.NumField(); j++ {
				//make name and Name to lowercase
				input.Name = strings.ToLower(input.Name)
				KeyName := strings.ToLower(permissionDetail.Type().Field(j).Name)
				if KeyName == input.Name {
					if !permissionDetail.Field(j).Interface().(bool) {
						return false, exception.UnauthorizedError{Message: "permission denied"}
					}
					return true, nil
				}
			}
		}
	}
	return false, exception.UnauthorizedError{Message: "permission denied"}
}
