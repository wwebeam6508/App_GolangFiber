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

	searchPipeline := bson.A{}
	if body.Search != "%%" && body.SearchFilter != "%%" {
		searchPipeline = append(searchPipeline, bson.M{body.SearchFilter: bson.M{"$regex": body.Search, "$options": "i"}})
	}
	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	input := model.GetCustomerInput{
		Page:         body.Page,
		PageSize:     body.PageSize,
		SortTitle:    body.SortTitle,
		SortType:     body.SortType,
		Search:       body.Search,
		SearchFilter: body.SearchFilter,
	}
	customerCountChan, errChan := make(chan int32), make(chan error)
	go service.GetCustomerCountService(searchPipelineGroup, customerCountChan, errChan)
	if len(errChan) > 0 {
		return exception.ErrorHandler(c, <-errChan)
	}
	customerChan, errChan := make(chan []model.GetCustomerResult), make(chan error)
	go func() {
		customer, err := service.GetCustomerService(input, searchPipelineGroup)
		if err != nil {
			errChan <- err
			return
		}
		customerChan <- customer
	}()
	if len(errChan) > 0 {
		return exception.ErrorHandler(c, <-errChan)
	}
	customerCount := <-customerCountChan
	customer := <-customerChan
	pages := common.PageArray(customerCount, input.PageSize, input.Page, 5)
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: commonentity.PaginationResponse{
			CurrentPage: input.Page,
			Pages:       pages,
			Data:        customer,
			LastPage:    int(math.Ceil(float64(customerCount) / float64(input.PageSize))),
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
	if body.Search == "" {
		body.Search = "%%"
	}
	if body.SearchFilter == "" {
		body.SearchFilter = "%%"
	}
	return body
}
