package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/location"
	service "PBD_backend_go/service/location"
	"math"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLocationController(c *fiber.Ctx) error {
	var query commonentity.PaginateInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	prePaginateCondition(&query)
	searchGroup, err := getSearchGroup(query.Search, query.SearchFilter)
	if err != nil {
		return err
	}

	//chan of allLocationCount
	allLocationCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetLocationCountService(searchGroup)
		if err != nil {
			errChan <- err
			allLocationCountChan <- 0
			return
		}
		allLocationCountChan <- int32(count)
		errChan <- nil
	}()
	//chan of wage
	resultChan := make(chan []model.GetLocationResult, 1)
	go func() {
		result, err := service.GetLocationService(query, searchGroup)
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
	allLocationCount := <-allLocationCountChan
	result := <-resultChan
	return c.Status(200).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(allLocationCount, query.PageSize, query.Page, 5),
		Data:        result,
		LastPage:    int(math.Ceil(float64(allLocationCount) / float64(query.PageSize))),
	})
}

func GetLocationByIDController(c *fiber.Ctx) error {
	input := model.GetLocationByIDInput{}
	if err := c.QueryParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	result, err := service.GetLocationByIDService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func UpdateLocationController(c *fiber.Ctx) error {
	query := model.UpdateLocationID{}
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(query)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	var input model.UpdateLocationInput
	if err := c.BodyParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := service.UpdateLocationService(input, query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func AddLocationController(c *fiber.Ctx) error {
	var input model.AddLocationInput
	if err := c.BodyParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	result, err := service.AddLocationService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func DeleteLocationController(c *fiber.Ctx) error {
	var input model.DeleteLocationID
	if err := c.QueryParser(&input); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := common.Validate(input)
	if validate != nil {
		return exception.ErrorHandler(c, validate)
	}
	err := service.DeleteLocationService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func getSearchGroup(search string, searchFilter string) (commonentity.SearchPipeline, error) {
	var searchFilterArray bson.A
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		searchFilterArray = bson.A{
			bson.M{
				searchFilter: bson.M{
					"$regex":   search,
					"$options": "i",
				},
			},
		}
	}
	return commonentity.SearchPipeline{
		Search:         search,
		SearchPipeline: searchFilterArray,
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
