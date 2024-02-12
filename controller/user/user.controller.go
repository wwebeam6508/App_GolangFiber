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
	getUserBodyCondition(&body)
	// searchPipeline as array
	searchPipeline, err := getSearchPipeline(body.Search, body.SearchFilter)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	// check is searchPipeline empty
	input := model.GetUserServiceInput{
		Page:      body.Page,
		PageSize:  body.PageSize,
		SortTitle: body.SortTitle,
		SortType:  body.SortType,
	}
	resultChan := make(chan []model.GetUserServiceResult, 1)
	errChan := make(chan error, 2)
	allUserCountChan := make(chan int32, 1)
	go func() {
		result, err := service.GetUserService(input, searchPipelineGroup)
		if err != nil {
			errChan <- err
			allUserCountChan <- 0
			return
		}
		resultChan <- result
		errChan <- nil
	}()
	go func() {
		allUserCount, err := service.GetAllUserCount(searchPipelineGroup)
		if err != nil {
			errChan <- err
			allUserCountChan <- 0
			return
		}
		allUserCountChan <- allUserCount
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	result := <-resultChan
	allUserCount := <-allUserCountChan

	pages := common.PageArray(allUserCount, input.PageSize, input.Page, 5)
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
	result, err := service.GetUserTypeService()
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
		return exception.ErrorHandler(c, exception.ValidationError{Message: "your rank has been disabled"})
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

func getSearchPipeline(search, searchFilter string) (bson.A, error) {
	searchPipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		if searchFilter == "userType" {
			searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "userType.name", Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}}}}})
		} else if searchFilter == "date" {
			split := strings.Split(search, ",")
			if len(split) != 2 {
				return searchPipeline, exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
				//time Parse
				dateSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return searchPipeline, exception.ValidationError{Message: "invalid date"}
				}
				searchPipeline = append(searchPipeline, bson.M{"createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(dateSearch)}})
			} else {
				dateStartSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return searchPipeline, exception.ValidationError{Message: "invalid date"}
				}
				dateEndSearch, err := time.Parse(time.RFC3339, split[1])
				if err != nil {
					return searchPipeline, exception.ValidationError{Message: "invalid date"}
				}
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateStartSearch)}, {Key: "$lte", Value: primitive.NewDateTimeFromTime(dateEndSearch)}}}}}})
			}
		} else {
			searchPipeline = append(searchPipeline, bson.M{searchFilter: bson.M{"$regex": search, "$options": "i"}})
		}
	}
	return searchPipeline, nil
}

func getUserBodyCondition(input *model.GetUserControllerInput) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.SortTitle == "" {
		input.SortTitle = "date"
	}
	if input.SortType == "" {
		input.SortType = "desc"
	}
}
