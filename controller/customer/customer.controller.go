package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/customer"
	service "PBD_backend_go/service/customer"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCustomerController(c *fiber.Ctx) error {
	var body model.GetCustomerInput
	//query not body
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	body = getCustomerBodyCondition(body)

	searchPipeline := getSearchPipeline(body.Search, body.SearchFilter)

	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	customerCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetCustomerCountService(searchPipelineGroup)
		if err != nil {
			errChan <- err
			customerCountChan <- 0
			return
		}
		customerCountChan <- count
		errChan <- nil
	}()
	err := <-errChan
	if err != nil {
		return exception.ErrorHandler(c, <-errChan)
	}
	customerChan := make(chan []model.GetCustomerResult, 1)
	go func() {
		customer, err := service.GetCustomerService(body, searchPipelineGroup)
		if err != nil {
			errChan <- err
			customerChan <- nil
			return
		}
		customerChan <- customer
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, <-errChan)
	}
	customerCount := <-customerCountChan
	customer := <-customerChan
	pages := common.PageArray(customerCount, body.PageSize, body.Page, 5)
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: commonentity.PaginationResponse{
			CurrentPage: body.Page,
			Pages:       pages,
			Data:        customer,
			LastPage:    int(math.Ceil(float64(customerCount) / float64(body.PageSize))),
		},
	})
}

func GetCustomerByIDController(c *fiber.Ctx) error {
	var body model.GetCustomerByIDInput
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validator := validator.New()
	err := validator.Struct(body)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}

	if body.CustomerID == "" {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "customerID is required"})
	}
	customer, err := service.GetCustomerByIDService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    customer,
	})
}

func AddCustomerController(c *fiber.Ctx) error {
	var body model.AddCustomerInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	id, err := service.AddCustomerService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    id,
	})

}

func UpdateCustomerController(c *fiber.Ctx) error {
	var query model.UpdateCustomerID
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	err := validate.Struct(query)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}

	var body model.UpdateCustomerInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err = service.UpdateCustomerService(body, query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func DeleteCustomerController(c *fiber.Ctx) error {
	var body model.DeleteCustomerInput
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	err = service.DeleteCustomerService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func getSearchPipeline(search, searchFilter string) bson.A {
	pipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}}}}})
	}
	return pipeline
}

func getCustomerBodyCondition(body model.GetCustomerInput) model.GetCustomerInput {
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
