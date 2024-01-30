package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/userManagement"
	jwtservice "PBD_backend_go/service/auth"
	service "PBD_backend_go/service/userManagement"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserController(c *fiber.Ctx) error {
	var body model.GetUserControllerInput
	//query not body
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	body = getUserBodyCondition(body)
	// searchPipeline as array
	searchPipeline := bson.A{}
	if body.Search != "%%" && body.SearchFilter != "%%" {
		// if searchFilter is "userType" then { "userType.name": { $regex: search, $options: "i" } }
		if body.SearchFilter == "userType" {
			searchPipeline = append(searchPipeline, bson.M{"userType.name": bson.M{"$regex": body.Search, "$options": "i"}})
		} else if body.SearchFilter == "date" {
			split := strings.Split(body.Search, ",")
			if len(split) != 2 {
				return exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
				//time Parse
				dateSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return exception.ErrorHandler(c, err)
				}
				searchPipeline = append(searchPipeline, bson.M{"createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(dateSearch)}})
			} else {
				dateStartSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return exception.ErrorHandler(c, err)
				}
				dateEndSearch, err := time.Parse(time.RFC3339, split[1])
				if err != nil {
					return exception.ErrorHandler(c, err)
				}
				searchPipeline = append(searchPipeline, bson.M{"createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(dateStartSearch), "$lte": primitive.NewDateTimeFromTime(dateEndSearch)}})
			}
		} else {
			searchPipeline = append(searchPipeline, bson.M{body.SearchFilter: bson.M{"$regex": body.Search, "$options": "i"}})
		}
	}
	// check is searchPipeline empty
	input := model.GetUserServiceInput{
		Page:           body.Page,
		PageSize:       body.PageSize,
		SortTitle:      body.SortTitle,
		SortType:       body.SortType,
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	result, err := service.GetUserService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func GetUserByIDController(c *fiber.Ctx) error {
	var body model.GetUserByIDInput
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	result, err := service.GetUserByIDService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func AddUserController(c *fiber.Ctx) error {
	var body model.AddUserInput
	if common.DenialIfSuperAdmin(body.UserTypeID) {
		return exception.ErrorHandler(c, exception.UnauthorizedError{Message: "cannot add super admin"})
	}
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := service.AddUserService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    nil,
	})
}

func UpdateUserController(c *fiber.Ctx) error {
	var body model.UpdateUserInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}

	// updateCondition
	if err := updateCondition(c, body); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := service.UpdateUserService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func DeleteUserController(c *fiber.Ctx) error {
	var body model.DeleteUserInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	// updateCondition
	if err := deleteCondition(c, body); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := service.DeleteUserService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func updateCondition(c *fiber.Ctx, input model.UpdateUserInput) error {
	//check empty
	if input.UserID == "" || input.SelfID == "" {
		return exception.ValidationError{Message: "invalid userID or selfID"}
	}
	//userID from authorization
	selfIDSplit := strings.Split(c.Get("Authorization"), " ")
	if len(selfIDSplit) != 2 {
		return exception.UnauthorizedError{Message: "permission denied"}
	}
	input.SelfID = selfIDSplit[1]
	//get userID from claim
	claims, err := jwtservice.VerifyJWT(input.SelfID)
	if err != nil {
		return err
	}
	input.SelfID = claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})["userID"].(string)
	//check if userID or selfID is empty

	if service.StopChangeItself(input.UserID, input.SelfID) {
		return exception.UnauthorizedError{Message: "cannot update itself"}
	}
	//get userTypeID from input.SelfID
	oldData, err := service.GetUserByIDService(model.GetUserByIDInput{UserID: input.UserID})
	if err != nil {
		return err
	}
	//check if userTypeID is super admin
	if service.StopChangeSuperAdmin(oldData.UserTypeID.Hex()) {
		return exception.UnauthorizedError{Message: "cannot update super admin"}
	}
	if service.StopChangeSuperAdmin(input.UserTypeID) {
		return exception.UnauthorizedError{Message: "cannot update super admin"}
	}

	return nil
}

func deleteCondition(c *fiber.Ctx, input model.DeleteUserInput) error {
	//check empty
	if input.UserID == "" {
		return exception.ValidationError{Message: "invalid userID"}
	}
	//get userID from authorization
	split := strings.Split(c.Get("Authorization"), " ")
	if len(split) != 2 {
		return exception.UnauthorizedError{Message: "permission denied"}
	}
	userID := split[1]
	//get userID from claim
	claims, err := jwtservice.VerifyJWT(userID)
	if err != nil {
		return err
	}
	userID = claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})["userID"].(string)
	//check if userID or selfID is empty
	if service.StopChangeItself(input.UserID, userID) {
		return exception.UnauthorizedError{Message: "cannot update itself"}
	}
	//get userTypeID from input.SelfID
	oldData, err := service.GetUserByIDService(model.GetUserByIDInput{UserID: input.UserID})
	if err != nil {
		return err
	}
	//check if userTypeID is super admin
	if service.StopChangeSuperAdmin(oldData.UserTypeID.Hex()) {
		return exception.UnauthorizedError{Message: "cannot update super admin"}
	}

	return nil
}

func getUserBodyCondition(input model.GetUserControllerInput) model.GetUserControllerInput {
	var result model.GetUserControllerInput
	if input.PageSize <= 0 {
		result.PageSize = 10
	}
	if input.SortTitle == "" {
		result.SortTitle = "date"
	}
	if input.SortType == "" {
		result.SortType = "desc"
	}
	if input.Search == "" {
		result.Search = "%%"
	} else {
		result.Search = "%" + input.Search + "%" // %input.Search%
	}
	if input.SearchFilter == "" {
		result.SearchFilter = "%%"
	} else {
		result.SearchFilter = "%" + input.SearchFilter + "%" // %body.SearchFilter%
	}
	return result
}
