package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/userType"
	authservice "PBD_backend_go/service/auth"
	service "PBD_backend_go/service/userType"
	"math"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserTypeController(c *fiber.Ctx) error {
	//get body
	var body model.GetUserTypeInput
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	body = getUserTypeBodyCondition(body)
	searchPipeline, err := getSearchPipeline(body.Search, body.SearchFilter)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	//get count
	allUserTypeCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetAllUserTypeCountService(searchPipelineGroup)
		if err != nil {
			errChan <- err
			allUserTypeCountChan <- 0
			return
		}
		allUserTypeCountChan <- count
		errChan <- nil
	}()
	//get userType
	resultChan := make(chan []model.GetUserTypeResult, 1)
	go func() {
		result, err := service.GetUserTypeService(body, searchPipelineGroup)
		if err != nil {
			errChan <- err
			resultChan <- nil
			return
		}
		resultChan <- result
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	result := <-resultChan
	allUserTypeCount := <-allUserTypeCountChan

	pages := common.PageArray(allUserTypeCount, body.PageSize, body.Page, 5)
	//filter rank
	result, err = filterRankGetUserTypeController(c, body, result)
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
			"lastPage":    math.Ceil(float64(allUserTypeCount) / float64(body.PageSize)),
		},
	})
}

func GetUserTypeByIDController(c *fiber.Ctx) error {
	//get id
	var query model.GetUserTypeByIDInput
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	//get userType
	result, err := service.GetUserTypeByIDService(query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func AddUserTypeController(c *fiber.Ctx) error {
	//get body
	var body model.AddUserTypeInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	//add userType
	result, err := service.AddUserTypeService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    result,
	})
}

func UpdateUserTypeController(c *fiber.Ctx) error {
	//get id
	var query model.UpdateUserTypeID
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	//get body
	var body model.UpdateUserTypeInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	//update userType
	err := service.UpdateUserTypeService(body, query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteUserTypeController(c *fiber.Ctx) error {
	//get body
	var body model.DeleteUserTypeID
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	//delete userType
	err := service.DeleteUserTypeService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func filterRankGetUserTypeController(c *fiber.Ctx, input model.GetUserTypeInput, result []model.GetUserTypeResult) ([]model.GetUserTypeResult, error) {
	//get rank from token
	token := c.Get("Authorization")
	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return nil, exception.ValidationError{Message: "invalid token"}
	}
	claims, err := authservice.VerifyJWT(splitToken[1])
	if err != nil {
		return nil, exception.ValidationError{Message: "invalid token"}
	}
	userData := claims.Claims.(jwt.MapClaims)["data"].(map[string]interface{})
	//check if userType is super admin
	rank := userData["userType"].(map[string]interface{})["rank"].(float64)
	rankInt32 := int32(rank)
	if rankInt32 == 0 {
		return nil, exception.ValidationError{Message: "your rank has been disabled"}
	}
	//check rank in result show only rank lower than user rank
	var resultFilter []model.GetUserTypeResult
	for _, v := range result {
		if v.Rank > rankInt32 {
			resultFilter = append(resultFilter, model.GetUserTypeResult{
				UserTypeID: v.UserTypeID,
				Name:       v.Name,
				Date:       v.Date,
				Rank:       v.Rank,
			})
		}
	}
	return resultFilter, nil
}

func getSearchPipeline(search string, searchFilter string) (bson.A, error) {
	searchPipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		if searchFilter == "date" {
			split := strings.Split(search, ",")
			if len(split) != 2 {
				return searchPipeline, exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
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
				searchPipeline = append(searchPipeline, bson.M{"createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(dateStartSearch), "$lte": primitive.NewDateTimeFromTime(dateEndSearch)}})
			}
		} else {
			searchPipeline = append(searchPipeline, bson.M{searchFilter: bson.M{"$regex": search, "$options": "i"}})
		}
	}
	return searchPipeline, nil
}

func getUserTypeBodyCondition(body model.GetUserTypeInput) model.GetUserTypeInput {
	if body.Page == 0 {
		body.Page = 1
	}
	if body.PageSize == 0 {
		body.PageSize = 10
	}
	if body.SortTitle == "" {
		body.SortTitle = "date"
	}
	if body.SortType == "" {
		body.SortType = "desc"
	}
	return body
}
