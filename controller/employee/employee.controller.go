package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/employee"
	service "PBD_backend_go/service/employee"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetEmployeeController(c *fiber.Ctx) error {
	var query model.GetEmployeeInput
	if err := c.QueryParser(&query); err != nil {
		return err
	}
	query = getEmployeeBodyCondition(query)
	searchPipeline, err := getSearchPipeline(query.Search, query.SearchFilter)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	searchPipelineGroup := commonentity.SearchPipeline{
		Search:         query.Search,
		SearchPipeline: searchPipeline,
	}
	employeeCountChan, errChan := make(chan int32, 1), make(chan error, 1)
	go func() {
		count, err := service.GetEmployeeCountService(searchPipelineGroup)
		if err != nil {
			errChan <- err
			employeeCountChan <- 0
			return
		}
		employeeCountChan <- count
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	employeeChan, errChan1 := make(chan []model.GetEmployeeResult, 1), make(chan error, 1)
	go func() {
		employee, err := service.GetEmployeeService(query, searchPipelineGroup)
		if err != nil {
			errChan1 <- err
			employeeChan <- nil
			return
		}
		employeeChan <- employee
		errChan1 <- nil
	}()
	err = <-errChan1
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	employeeCount := <-employeeCountChan
	employee := <-employeeChan
	return c.Status(fiber.StatusOK).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		Data:        employee,
		CurrentPage: query.Page,
		LastPage:    int(math.Ceil(float64(employeeCount) / float64(query.PageSize))),
		Pages:       common.PageArray(employeeCount, query.PageSize, query.Page, 5),
	})

}

func getSearchPipeline(search, searchFilter string) (bson.A, error) {
	pipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		if searchFilter == "bornDate" || searchFilter == "joinedDate" {
			split := strings.Split(search, ",")
			if len(split) != 2 {
				return pipeline, exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
				//time Parse
				dateSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return pipeline, exception.ValidationError{Message: "invalid date"}
				}
				pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateSearch)}}}}}})
			} else {
				dateStartSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return pipeline, exception.ValidationError{Message: "invalid date"}
				}
				dateEndSearch, err := time.Parse(time.RFC3339, split[1])
				if err != nil {
					return pipeline, exception.ValidationError{Message: "invalid date"}
				}
				pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateStartSearch)}, {Key: "$lte", Value: primitive.NewDateTimeFromTime(dateEndSearch)}}}}}})
			}
		} else if searchFilter == "salary" {
			pipeline = append(pipeline,
				bson.D{
					{Key: "$match", Value: bson.D{
						{Key: searchFilter, Value: bson.D{
							{Key: "$gte", Value: strings.Split(search, ",")[0]},
							{Key: "$lte", Value: strings.Split(search, ",")[1]},
						}},
					}},
				},
			)
		} else {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}}}}})
		}
	}
	return pipeline, nil
}

func getEmployeeBodyCondition(body model.GetEmployeeInput) model.GetEmployeeInput {
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
