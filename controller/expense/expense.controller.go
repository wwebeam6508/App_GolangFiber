package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/expense"
	service "PBD_backend_go/service/expense"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetExpenseController(c *fiber.Ctx) error {
	query := model.GetExpenseInput{}
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	query = getExpenseBodyCondition(query)
	searchPipeline, err := getSearchPipeline(query.Search, query.SearchFilter)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	searchPipelineGroup := model.SearchPipeline{
		Search:         query.Search,
		SearchPipeline: searchPipeline,
	}
	expenseCountChan, errChan := make(chan int32, 1), make(chan error, 1)
	go func() {
		count, err := service.GetExpenseCountService(searchPipelineGroup)
		if err != nil {
			errChan <- err
			expenseCountChan <- 0
			return
		}
		expenseCountChan <- count
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	expenseChan, errChan1 := make(chan []model.GetExpenseServiceResult, 1), make(chan error, 1)
	go func() {
		expense, err := service.GetExpenseService(query, searchPipelineGroup)
		if err != nil {
			errChan1 <- err
			expenseChan <- nil
			return
		}
		expenseChan <- expense
		errChan1 <- nil
	}()
	err1 := <-errChan1
	if err1 != nil {
		return exception.ErrorHandler(c, err1)
	}
	expenseCount := <-expenseCountChan
	expense := <-expenseChan
	return c.Status(fiber.StatusOK).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: query.Page,
		Pages:       common.PageArray(expenseCount, query.PageSize, query.Page, 5),
		Data:        expense,
		LastPage:    int(math.Ceil(float64(expenseCount) / float64(query.PageSize))),
	})
}

func GetExpenseByIDController(c *fiber.Ctx) error {
	query := model.GetExpenseByIDInput{}
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	expense, err := service.GetExpenseByIDService(query)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    expense,
	})
}

func AddExpenseController(c *fiber.Ctx) error {
	body := model.AddExpenseInput{}
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	id, err := service.AddExpenseService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    id,
	})
}

func UpdateExpenseController(c *fiber.Ctx) error {
	query := model.UpdateExpenseID{}
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(query); err != nil {
		return exception.ErrorHandler(c, err)
	}

	body := model.UpdateExpenseInput{}
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	if err := validate.Struct(body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := service.UpdateExpenseService(body, query.ExpenseID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func DeleteExpenseController(c *fiber.Ctx) error {
	body := model.DeleteExpenseInput{}
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	err := service.DeleleExpenseService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func GetWorkTitleController(c *fiber.Ctx) error {
	workTitle, err := service.GetWorkTitleService()
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    workTitle,
	})
}

func GetCustomerNameController(c *fiber.Ctx) error {
	customerName, err := service.GetCustomerNameService()
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(200).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    customerName,
	})
}

func getSearchPipeline(search string, searchFilter string) (bson.A, error) {
	pipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		if searchFilter == "work" {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "workRef.title", Value: bson.D{
				{Key: "$regex", Value: search},
			}}}}})
		} else if searchFilter == "expense" {
			//add field that sum of lists.price
			pipeline = append(pipeline,
				bson.M{
					"$addFields": bson.M{
						"listSum": bson.M{
							"$reduce": bson.M{
								"input":        "$lists",
								"initialValue": 0,
								"in": bson.M{
									"$add": bson.A{"$$value", "$$this.price"},
								},
							},
						},
					},
				},
			)
			split := strings.Split(search, ",")
			startExpense, _ := strconv.ParseFloat(split[0], 64)
			endExpense, _ := strconv.ParseFloat(split[1], 64)
			if startExpense > endExpense {
				return pipeline, exception.ValidationError{Message: "invalid expense"}
			}
			if common.IsEmpty(startExpense) {
				pipeline = append(pipeline, bson.M{"$match": bson.M{"listSum": bson.M{"$lte": endExpense}}})
			} else if common.IsEmpty(endExpense) {
				pipeline = append(pipeline, bson.M{"$match": bson.M{"listSum": bson.M{"$gte": startExpense}}})
			} else {
				pipeline = append(pipeline,
					bson.D{
						{Key: "$match", Value: bson.D{
							{Key: "listSum", Value: bson.D{
								{Key: "$gte", Value: startExpense},
								{Key: "$lte", Value: endExpense},
							}},
						}},
					},
				)
			}
		} else if searchFilter == "date" {
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
				pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "createdAt", Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateStartSearch)}, {Key: "$lte", Value: primitive.NewDateTimeFromTime(dateEndSearch)}}}}}})
			}
		} else {
			pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$regex", Value: search}}}}}})
		}
	}
	return pipeline, nil
}

func getExpenseBodyCondition(query model.GetExpenseInput) model.GetExpenseInput {
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
	return query
}
