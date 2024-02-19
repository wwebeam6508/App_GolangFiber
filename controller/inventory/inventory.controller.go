package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/inventory"
	service "PBD_backend_go/service/inventory"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetInventoryController(c *fiber.Ctx) error {
	var query commonentity.PaginateInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	prePaginateCondition(&query)
	searchGroup := getSearchGroup(query.Search, query.SearchFilter)
	//chan of allLocationCount
	allInventoryCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetInventoryCountService(searchGroup)
		if err != nil {
			errChan <- err
			allInventoryCountChan <- 0
			return
		}
		allInventoryCountChan <- int32(count)
		errChan <- nil
	}()
	//chan of inventoryType
	resultChan := make(chan []model.GetInventoryResult, 1)
	go func() {
		result, err := service.GetInventoryService(query, searchGroup)
		if err != nil {
			errChan <- err
			resultChan <- nil
			return
		}
		resultChan <- result
		errChan <- nil
	}()
	//get all chan
	allInventoryCount, result, err := <-allInventoryCountChan, <-resultChan, <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(allInventoryCount, query.PageSize, query.Page, 5),
		Data:        result,
		LastPage:    int(math.Ceil(float64(allInventoryCount) / float64(query.PageSize))),
	})
}

func GetInventoryByIDController(c *fiber.Ctx) error {
	var input model.GetInventoryByIDInput
	if err := c.QueryParser(&input); err != nil {
		return err
	}
	result, err := service.GetInventoryByID(input)
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

func AddInventoryController(c *fiber.Ctx) error {
	var input model.AddInventoryInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	res, err := service.AddInventoryService(input)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(201).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusCreated,
		Message: "Inventory Type Created",
		Data:    res,
	})
}

func UpdateInventoryController(c *fiber.Ctx) error {
	var ID model.UpdateInventoryID
	if err := c.QueryParser(&ID); err != nil {
		return err
	}
	if err := common.Validate(ID); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	var input model.UpdateInventoryInput
	if err := c.BodyParser(&input); err != nil {
		return err
	}
	err := service.UpdateInventoryService(input, ID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteInventoryController(c *fiber.Ctx) error {
	var ID model.DeleteInventoryID
	if err := c.QueryParser(&ID); err != nil {
		return err
	}
	if err := common.Validate(ID); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	err := service.DeleteInventoryService(ID)
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
		if searchFilter == "price" {
			//split search by comma
			searchArr := strings.Split(search, ",")
			gte, _ := strconv.ParseFloat(searchArr[0], 64)
			lte, _ := strconv.ParseFloat(searchArr[1], 64)
			pricePipeline := bson.A{}
			if !common.IsEmpty(gte) {
				pricePipeline = append(pricePipeline, bson.M{"price": bson.M{"$gte": gte}})
			}
			if !common.IsEmpty(lte) {
				pricePipeline = append(pricePipeline, bson.M{"price": bson.M{"$lte": lte}})
			}
			searchPipeline = append(searchPipeline, bson.M{"$and": pricePipeline})
		} else if searchFilter == "quantity" {
			//split search by comma
			searchArr := strings.Split(search, ",")
			gte, _ := strconv.ParseFloat(searchArr[0], 64)
			lte, _ := strconv.ParseFloat(searchArr[1], 64)
			quantityPipeline := bson.A{}
			if !common.IsEmpty(gte) {
				quantityPipeline = append(quantityPipeline, bson.M{"quantity": bson.M{"$gte": gte}})
			}
			if !common.IsEmpty(lte) {
				quantityPipeline = append(quantityPipeline, bson.M{"quantity": bson.M{"$lte": lte}})
			}
			searchPipeline = append(searchPipeline, bson.M{"$and": quantityPipeline})
		} else if searchFilter == "inventoryType" {
			searchPipeline = bson.A{
				bson.M{
					"inventoryTypeDetail.name": bson.M{
						"$regex":   search,
						"$options": "i",
					},
				},
			}
		} else {
			searchPipeline = bson.A{
				bson.M{
					searchFilter: bson.M{
						"$regex":   search,
						"$options": "i",
					},
				},
			}
		}
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
