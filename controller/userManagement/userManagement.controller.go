package controller

import (
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/userManagement"
	service "PBD_backend_go/service/userManagement"
	"strings"
	"time"

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

func getUserBodyCondition(input model.GetUserControllerInput) model.GetUserControllerInput {
	var result model.GetUserControllerInput
	if pageSize := input.PageSize; pageSize <= 0 {
		result.PageSize = 10
	}
	if sortTitle := input.SortTitle; sortTitle == "" {
		result.SortTitle = "date"
	}
	if sortType := input.SortType; sortType == "" {
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
