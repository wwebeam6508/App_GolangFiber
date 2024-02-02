package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/user"
	authservice "PBD_backend_go/service/auth"
	service "PBD_backend_go/service/user"
	"math"
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
	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
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
	resultChan := make(chan []model.GetUserServiceResult)
	errChan := make(chan error)
	allUserCountChan := make(chan int32)
	go func() {
		result, err := service.GetUserService(input, searchPipelineGroup)
		resultChan <- result
		errChan <- err
	}()
	go func() {
		allUserCount := service.GetAllUserCount(searchPipelineGroup)
		allUserCountChan <- allUserCount
	}()
	result := <-resultChan
	err := <-errChan
	allUserCount := <-allUserCountChan

	pages := common.PageArray(allUserCount, input.PageSize, input.Page, 5)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: bson.M{
			"currentPage": body.Page,
			"pages":       pages,
			"data":        result,
			"lastPage":    math.Ceil(float64(allUserCount) / float64(body.PageSize)),
		},
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
	result, err := service.AddUserService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    result,
	})
}

func UpdateUserController(c *fiber.Ctx) error {
	var query model.UpdateUserID
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}

	var body model.UpdateUserInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}

	err := service.UpdateUserService(body, query)
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

func GetUserTypeNameController(c *fiber.Ctx) error {
	resultChan := make(chan []model.GetUserTypeServiceResult)
	errChan := make(chan error)
	go func() {
		result, err := service.GetUserTypeService()
		resultChan <- result
		errChan <- err
	}()
	result := <-resultChan
	err := <-errChan

	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	//get rank from token
	token := c.Get("Authorization")
	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "invalid token"})
	}
	claims, err := authservice.VerifyJWT(splitToken[1])
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	userData := claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	//check if userType is super admin
	rank := userData["userType"].(map[string]interface{})["rank"].(float64)
	rankInt32 := int32(rank)
	if rankInt32 == 0 {
		return exception.ErrorHandler(c, exception.UnauthorizedError{Message: "your rank has been disabled"})
	}
	//check rank in result show only rank lower than user rank
	var resultFilter []model.GetUserTypeNameResult
	for _, v := range result {
		if v.Rank > rankInt32 {
			resultFilter = append(resultFilter, model.GetUserTypeNameResult{
				ID:   v.ID,
				Name: v.Name,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    resultFilter,
	})
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
