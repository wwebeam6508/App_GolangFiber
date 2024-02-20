package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/requisition"
	service "PBD_backend_go/service/requisition"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetRequisitionController(c *fiber.Ctx) error {
	var query commonentity.PaginateInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	prePaginateCondition(&query)
	searchGroup, err := getSearchGroup(query.Search, query.SearchFilter)
	if err != nil {
		return err
	}
	//chan of allRequisitionCount
	allRequisitionCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetRequisitionCountService(searchGroup)
		if err != nil {
			errChan <- err
			allRequisitionCountChan <- 0
			return
		}
		allRequisitionCountChan <- count
		errChan <- nil
	}()
	// chan of requisition
	requisitionChan := make(chan []model.GetRequisitionResult, 1)
	go func() {
		requisition, err := service.GetRequisitionService(query, searchGroup)
		if err != nil {
			errChan <- err
			requisitionChan <- nil
			return
		}
		requisitionChan <- requisition
		errChan <- nil
	}()
	//get all chan
	allRequisitionCount, result, err := <-allRequisitionCountChan, <-requisitionChan, <-errChan
	if err != nil {
		return err
	}
	return c.Status(200).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(allRequisitionCount, query.PageSize, query.Page, 5),
		Data:        result,
		LastPage:    int(math.Ceil(float64(allRequisitionCount) / float64(query.PageSize))),
	})
}

func GetRequisitionByIDController(c *fiber.Ctx) error {
	var input model.GetRequisitionByIDInput
	//string to objectID
	if err := c.QueryParser(&input); err != nil {
		return err
	}
	//validate
	if err := common.Validate(input); err != nil {
		return exception.ValidationError{Message: err.Error()}
	}
	if common.IsEmpty(input.RequisitionID) {
		return exception.ValidationError{Message: "RequisitionID is required"}
	}
	result, err := service.GetRequisitionByIDService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func AddRequisitionController(c *fiber.Ctx) error {
	var input model.AddRequisitionInput
	//parse body
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	//validate
	if err := common.Validate(input); err != nil {
		return exception.ValidationError{Message: err.Error()}
	}
	//add
	result, err := service.AddRequisitionService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(201).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Success",
		Data:    result,
	})
}

func UpdateRequisitionStatusService(c *fiber.Ctx) error {
	var query model.UpdateRequisitionStatusID
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	//validate
	if err := common.Validate(query); err != nil {
		return exception.ValidationError{Message: err.Error()}
	}
	var input model.UpdateRequisitionStatusInput
	//parse body
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	//validate
	if err := common.Validate(input); err != nil {
		return exception.ValidationError{Message: err.Error()}
	}
	//update
	err := service.UpdateRequisitionStatusService(query, input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteRequisitionController(c *fiber.Ctx) error {
	var query model.DeleteRequisitionID
	//parse body
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	//validate
	if err := common.Validate(query); err != nil {
		return exception.ValidationError{Message: err.Error()}
	}
	//delete
	err := service.DeleteRequisitionService(query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func getSearchGroup(search string, searchFilter string) (commonentity.SearchPipeline, error) {
	searchPipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
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
		} else if searchFilter == "employee" {
			//search employeeDetail.FirstName and employeeDetail.LastName already has lookup
			searchName := strings.Split(search, " ")
			//check first and last is not empty
			searchNamePipeline := bson.A{}
			if !common.IsEmpty(searchName[0]) {
				searchNamePipeline = append(searchNamePipeline, bson.M{"employeeDetail.FirstName": bson.M{"$regex": searchName[0], "$options": "i"}})
			}
			if !common.IsEmpty(searchName[1]) {
				searchNamePipeline = append(searchNamePipeline, bson.M{"employeeDetail.LastName": bson.M{"$regex": searchName[1], "$options": "i"}})
			}
			searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{"$and": searchNamePipeline}})

		} else if searchFilter == "inventory" {
			//split search by space
			searchInventory := strings.Split(search, " ")
			searchInventoryPipeline := bson.A{}
			for _, v := range searchInventory {
				//search inventoryDetail.Name that is array
				searchInventoryPipeline = append(searchInventoryPipeline, bson.M{"inventoryDetail.Name": bson.M{"$regex": v, "$options": "i"}})
			}
			searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{"$and": searchInventoryPipeline}})
		} else {
			searchPipeline = append(searchPipeline, bson.M{"$match": bson.M{searchFilter: bson.M{"$regex": search, "$options": "i"}}})
		}
	}
	return commonentity.SearchPipeline{
		Search:         search,
		SearchPipeline: searchPipeline,
	}, nil
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
