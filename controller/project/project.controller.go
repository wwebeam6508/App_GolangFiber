package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/project"
	service "PBD_backend_go/service/project"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProjectController(c *fiber.Ctx) error {
	var body model.GetProjectInput
	//query not body
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	body = getProjectBodyCondition(body)
	// searchPipeline as array
	searchPipeline, err := getSearchPipeline(body.Search, body.SearchFilter)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	searchPipelineGroup := model.SearchPipeline{
		Search:         body.Search,
		SearchPipeline: searchPipeline,
	}
	projectCountChan, errChan := make(chan int32, 1), make(chan error, 2)
	go func() {
		count, err := service.GetProjectCountService(searchPipelineGroup)
		if err != nil {
			errChan <- err
			projectCountChan <- 0
			return
		}
		projectCountChan <- count
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	projectChan := make(chan []model.GetProjectServiceResult, 1)
	go func() {
		project, err := service.GetProjectService(body, searchPipelineGroup)
		if err != nil {
			errChan <- err
			projectChan <- nil
			return
		}
		projectChan <- project
		errChan <- nil
	}()
	err = <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	projectCount := <-projectCountChan
	project := <-projectChan
	pages := common.PageArray(projectCount, body.PageSize, body.Page, 5)
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: commonentity.PaginationResponse{
			CurrentPage: body.Page,
			Pages:       pages,
			Data:        project,
			LastPage:    int(math.Ceil(float64(projectCount) / float64(body.PageSize))),
		},
	})
}

func getSearchPipeline(search string, searchFilter string) (bson.A, error) {
	searchPipeline := bson.A{}
	if search != "%%" && searchFilter != "%%" {
		// if searchFilter is "customer" then { "customer.name": { $regex: search, $options: "i" } }
		if searchFilter == "customer" {
			searchPipeline = append(searchPipeline, bson.M{"customer.name": bson.M{"$regex": search, "$options": "i"}})
		} else if searchFilter == "date" || searchFilter == "dateEnd" {
			split := strings.Split(search, ",")
			if len(split) != 2 {
				return searchPipeline, exception.ValidationError{Message: "invalid date"}
			}
			if split[1] == "" {
				dateSearch, err := time.Parse(time.RFC3339, split[0])
				if err != nil {
					return searchPipeline, exception.ValidationError{Message: "invalid date"}
				}
				searchPipeline = append(searchPipeline, bson.M{searchFilter: bson.M{"$gte": primitive.NewDateTimeFromTime(dateSearch)}})
			} else {
				dateStartSearch, _ := time.Parse(time.RFC3339, split[0])
				dateEndSearch, _ := time.Parse(time.RFC3339, split[1])
				searchPipeline = append(searchPipeline, bson.M{searchFilter: bson.M{"$gte": primitive.NewDateTimeFromTime(dateStartSearch), "$lte": primitive.NewDateTimeFromTime(dateEndSearch)}})
			}
		} else {
			searchPipeline = append(searchPipeline, bson.M{searchFilter: bson.M{"$regex": search, "$options": "i"}})
		}
	}
	return searchPipeline, nil
}

func getProjectBodyCondition(body model.GetProjectInput) model.GetProjectInput {
	if body.Page <= 0 {
		body.Page = 1
	}
	if body.PageSize <= 0 {
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
