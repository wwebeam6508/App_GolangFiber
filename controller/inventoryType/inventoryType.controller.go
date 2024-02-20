package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/inventoryType"
	service "PBD_backend_go/service/inventoryType"
	"math"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetInventoryTypeController(c *fiber.Ctx) error {
	var query commonentity.PaginateInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	prePaginateCondition(&query)
	searchGroup := getSearchGroup(query.Search, query.SearchFilter)
	//chan of allLocationCount
	allInventoryTypeCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetInventoryTypeCountService(searchGroup)
		if err != nil {
			errChan <- err
			allInventoryTypeCountChan <- 0
			return
		}
		allInventoryTypeCountChan <- int32(count)
		errChan <- nil
	}()
	//chan of inventoryType
	resultChan := make(chan []model.GetInventoryTypeResult, 1)
	go func() {
		result, err := service.GetInventoryTypeService(query, searchGroup)
		if err != nil {
			errChan <- err
			resultChan <- nil
			return
		}
		resultChan <- result
		errChan <- nil
	}()
	//get all chan
	allInventoryTypeCount, result, err := <-allInventoryTypeCountChan, <-resultChan, <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(allInventoryTypeCount, query.PageSize, query.Page, 5),
		Data:        result,
		LastPage:    int(math.Ceil(float64(allInventoryTypeCount) / float64(query.PageSize))),
	})
}

func GetInventoryTypeByIDController(c *fiber.Ctx) error {
	var input model.GetInventoryByIDInput
	if err := c.QueryParser(&input); err != nil {
		return err
	}
	result, err := service.GetInventoryTypeByID(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	if result == nil {
		return exception.ErrorHandler(c, exception.NotFoundError{Message: "Not Found"})
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func AddInventoryTypeController(c *fiber.Ctx) error {
	var input model.AddInventoryTypeInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	res, err := service.AddInventoryTypeService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(201).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Inventory Type Created",
		Data:    res,
	})
}

func UpdateInventoryTypeController(c *fiber.Ctx) error {
	var ID model.UpdateInventoryTypeID
	if err := c.QueryParser(&ID); err != nil {
		return err
	}
	if err := common.Validate(ID); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	var input model.UpdateInventoryTypeInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	err := service.UpdateInventoryTypeService(input, ID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteInventoryTypeController(c *fiber.Ctx) error {
	var ID model.DeleteInventoryTypeID
	if err := c.QueryParser(&ID); err != nil {
		return err
	}
	if err := common.Validate(ID); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	err := service.DeleteInventoryTypeService(ID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func getSearchGroup(search string, searchFilter string) commonentity.SearchPipeline {
	searchPipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$regex", Value: search}}}}}})
	}
	return commonentity.SearchPipeline{
		Search:         search,
		SearchPipeline: searchPipeline,
	}
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
