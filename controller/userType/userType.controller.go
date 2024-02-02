package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/userType"
	authservice "PBD_backend_go/service/auth"
	service "PBD_backend_go/service/userType"
	"fmt"
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
	searchPipeline := bson.A{}
	if body.Search != "%%" && body.SearchFilter != "%%" {
		if body.SearchFilter == "date" {
			split := strings.Split(body.Search, ",")
			if len(split) != 2 {
				return exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
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
	input := model.GetUserTypeInput{
		Page:         body.Page,
		PageSize:     body.PageSize,
		SortTitle:    body.SortTitle,
		SortType:     body.SortType,
		Search:       body.Search,
		SearchFilter: body.SearchFilter,
	}
	allUserTypeCount := service.GetAllUserTypeCountService(searchPipelineGroup)
	result, err := service.GetUserTypeService(input, searchPipelineGroup)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	pages := common.PageArray(allUserTypeCount, input.PageSize, input.Page, 5)
	//filter rank
	result, err = filterRankGetUserTypeController(c, input, result)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: bson.M{
			"currentPage": input.Page,
			"pages":       pages,
			"data":        result,
			"lastPage":    math.Ceil(float64(allUserTypeCount) / float64(input.PageSize)),
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
		return nil, exception.UnauthorizedError{Message: "your rank has been disabled"}
	}
	//check rank in result show only rank lower than user rank
	var resultFilter []model.GetUserTypeResult
	for _, v := range result {
		if v.Rank > rankInt32 {
			fmt.Println(v.Rank, rankInt32)
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
	if body.Search == "" {
		body.Search = "%%"
	}
	if body.SearchFilter == "" {
		body.SearchFilter = "%%"
	}
	return body
}
