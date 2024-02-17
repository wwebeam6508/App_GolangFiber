package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/wage"
	service "PBD_backend_go/service/wage"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetWageController(c *fiber.Ctx) error {
	var query commonentity.PaginateInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	prePaginateCondition(&query)
	searchGroup, err := getSearchGroup(query.Search, query.SearchFilter)
	if err != nil {
		return err
	}

	//chan of allWageCount
	allWageCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetWageCountService(searchGroup)
		if err != nil {
			errChan <- err
			allWageCountChan <- 0
			return
		}
		allWageCountChan <- int32(count)
		errChan <- nil
	}()
	//chan of wage
	resultChan := make(chan []model.GetWageResult, 1)
	go func() {
		result, err := service.GetWageService(query, searchGroup)
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
	allWageCount := <-allWageCountChan
	result := <-resultChan
	return c.Status(200).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(allWageCount, query.PageSize, query.Page, 5),
		Data:        result,
		LastPage:    int(math.Ceil(float64(allWageCount) / float64(query.PageSize))),
	})
}

func GetWageByIDController(c *fiber.Ctx) error {
	var input model.GetWageByIDInput
	if err := c.QueryParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	result, err := service.GetWageByIDService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func AddWageController(c *fiber.Ctx) error {
	var input model.AddWageInput
	if err := c.BodyParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	result, err := service.AddWageService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func UpdateWageController(c *fiber.Ctx) error {
	var inputID model.UpdateWageID
	if err := c.QueryParser(&inputID); err != nil {
		return exception.ErrorHandler(c, err)
	}

	var input model.UpdateWageInput
	if err := c.BodyParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	err := common.Validate(inputID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if err := service.UpdateWageService(inputID, input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteWageController(c *fiber.Ctx) error {
	var input model.DeleteWageInput
	if err := c.QueryParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	err := service.DeleteWageService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func getSearchGroup(search string, searchFilter string) (commonentity.SearchPipeline, error) {
	searchPipeline := bson.A{}
	if search != "" && searchFilter != "" {
		if searchFilter == "date" {
			split := strings.Split(search, ",")
			if len(split) != 2 {
				return commonentity.SearchPipeline{}, exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
				//time Parse
				dateSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return commonentity.SearchPipeline{}, exception.ValidationError{Message: "invalid date"}
				}
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateSearch)}}}}}})
			} else {
				dateStartSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return commonentity.SearchPipeline{}, exception.ValidationError{Message: "invalid date"}
				}
				dateEndSearch, err := time.Parse(time.RFC3339, split[1])
				if err != nil {
					return commonentity.SearchPipeline{}, exception.ValidationError{Message: "invalid date"}
				}
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateStartSearch)}, {Key: "$lte", Value: primitive.NewDateTimeFromTime(dateEndSearch)}}}}}})
			}
		} else if searchFilter == "allWages" {
			addField := bson.M{
				"$addFields": bson.M{
					"allWages": bson.M{
						"$reduce": bson.M{
							"input":        "$employee",
							"initialValue": 0,
							"in": bson.M{
								"$add": bson.A{"$$value", "$$this.wage"},
							},
						},
					},
				},
			}
			searchPipeline = append(searchPipeline, addField)
			split := strings.Split(search, ",")
			startWage, _ := strconv.ParseFloat(split[0], 64)
			endWage, _ := strconv.ParseFloat(split[1], 64)
			if startWage > endWage {
				return commonentity.SearchPipeline{}, exception.ValidationError{Message: "invalid Wage"}
			}
			if common.IsEmpty(startWage) {
				searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{searchFilter: bson.M{"$lte": endWage}}})
			} else if common.IsEmpty(endWage) {
				searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{searchFilter: bson.M{"$gte": startWage}}})
			} else {
				searchPipeline = append(searchPipeline,
					bson.D{
						{Key: "$match", Value: bson.D{
							{Key: searchFilter, Value: bson.D{
								{Key: "$gte", Value: startWage},
								{Key: "$lte", Value: endWage},
							}},
						}},
					},
				)
			}
		} else if searchFilter == "employee_name" {
			lookupStage := bson.M{"$lookup": bson.M{"from": "employee", "localField": "employee.employeeID", "foreignField": "employeeID", "as": "employees"}}
			unwindStage := bson.M{"$unwind": bson.M{"path": "$employees", "preserveNullAndEmptyArrays": true}}
			searchPipeline = append(searchPipeline, lookupStage, unwindStage)
			searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{"employees.name": bson.M{"$regex": search, "$options": "i"}}})
		} else {
			searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{searchFilter: bson.M{"$regex": search, "$options": "i"}}})
		}
	}
	return commonentity.SearchPipeline{Search: search, SearchPipeline: searchPipeline}, nil
}

func prePaginateCondition(query *commonentity.PaginateInput) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}
	if query.SortTitle == "" {
		query.SortTitle = "date"
	}
	if query.SortType == "" {
		query.SortType = "desc"
	}
}
