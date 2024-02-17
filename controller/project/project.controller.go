package controller

import (
	"PBD_backend_go/common"
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/project"
	service "PBD_backend_go/service/project"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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
	return c.Status(fiber.StatusOK).JSON(commonentity.PaginationResponse{
		Code:        fiber.StatusOK,
		Message:     "Success",
		CurrentPage: body.Page,
		Pages:       common.PageArray(projectCount, body.PageSize, body.Page, 5),
		Data:        project,
		LastPage:    int(math.Ceil(float64(projectCount) / float64(body.PageSize))),
	})
}

func GetProjectByIDController(c *fiber.Ctx) error {
	var body model.GetProjectByIDInput
	if err := c.QueryParser(&body); err != nil {
		return exception.ErrorHandler(c, err)
	}
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}

	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "invalid project id"})
	}
	project, err := service.GetProjectByIDService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    project,
	})

}

func AddProjectController(c *fiber.Ctx) error {
	var body model.AddProjectInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(body)

	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	if body.DateEnd.Before(body.Date) {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "date end must be after date"})
	}
	projectID, err := service.AddProjectService(body)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}

	if len(body.Images) > 0 {
		ImagesUrl := make([]string, len(body.Images))
		for i, image := range body.Images {
			imageIndex := strconv.Itoa(i)
			url, err := common.UploadImageToStorage("works", projectID.Hex()+"_"+imageIndex, image)
			if err != nil {
				service.DeleteProjectService(projectID.Hex())
				for i := 0; i < len(body.Images); i++ {
					common.DeleteImageFromStorage("works", projectID.Hex()+"_"+strconv.Itoa(i))
				}
				return exception.ErrorHandler(c, err)
			}
			ImagesUrl[i] = url
		}
		err = service.UpdateProjectService(model.UpdateProjectInput{Images: ImagesUrl}, projectID.Hex())
		if err != nil {
			return exception.ErrorHandler(c, err)
		}
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    projectID,
	})
}

func UpdateProjectController(c *fiber.Ctx) error {
	var query model.UpdateProjectID
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	var body model.UpdateProjectInput
	if err := c.BodyParser(&body); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(body)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	err = validate.Struct(query)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	if body.DateEnd.Before(body.Date) {
		return exception.ErrorHandler(c, exception.ValidationError{Message: "date end must be after date"})
	}
	err = service.UpdateProjectService(body, query.ProjectID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	if len(body.Images) > 0 {
		ImagesUrl := make([]string, len(body.Images))
		for i, image := range body.Images {
			imageIndex := strconv.Itoa(i)
			url, err := common.UploadImageToStorage("project", query.ProjectID+"_"+imageIndex, image)
			if err != nil {
				return exception.ErrorHandler(c, err)
			}
			ImagesUrl[i] = url
		}
		err = service.UpdateProjectService(model.UpdateProjectInput{Images: ImagesUrl}, query.ProjectID)
		if err != nil {
			return exception.ErrorHandler(c, err)
		}

	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func DeleteProjectController(c *fiber.Ctx) error {
	var query model.DeleteProjectInput
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	validate := validator.New()
	err := validate.Struct(query)
	if err != nil {
		return exception.ErrorHandler(c, exception.ValidationError{Message: err.Error()})
	}
	err = service.DeleteProjectService(query.ProjectID)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func GetCustomerNameController(c *fiber.Ctx) error {

	result, err := service.GetCustomerNameService()
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    result,
	})
}

func getSearchPipeline(search string, searchFilter string) (bson.A, error) {
	searchPipeline := bson.A{}
	if !common.IsEmpty(search) && !common.IsEmpty(searchFilter) {
		if searchFilter == "customer" {
			//search customer name
			searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "customer.name", Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}}}}})
		} else if searchFilter == "profit" {
			// split search by comma
			split := strings.Split(search, ",")
			// convert to float
			profitStart, _ := strconv.ParseFloat(split[0], 64)
			profitEnd, _ := strconv.ParseFloat(split[1], 64)
			if profitStart > profitEnd {
				return searchPipeline, exception.ValidationError{Message: "invalid profit"}
			}
			if common.IsEmpty(profitStart) {
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "profit", Value: bson.D{{Key: "$lte", Value: profitEnd}}}}}})
			} else if common.IsEmpty(profitEnd) {
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "profit", Value: bson.D{{Key: "$gte", Value: profitStart}}}}}})
			} else {
				// add to pipeline
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "profit", Value: bson.D{{Key: "$gte", Value: profitStart}, {Key: "$lte", Value: profitEnd}}}}}})
			}

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
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateSearch)}}}}}})
			} else {
				dateStartSearch, _ := time.Parse(time.RFC3339, split[0])
				dateEndSearch, _ := time.Parse(time.RFC3339, split[1])
				searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$gte", Value: primitive.NewDateTimeFromTime(dateStartSearch)}, {Key: "$lte", Value: primitive.NewDateTimeFromTime(dateEndSearch)}}}}}})
			}
		} else {
			searchPipeline = append(searchPipeline, bson.D{{Key: "$match", Value: bson.D{{Key: searchFilter, Value: bson.D{{Key: "$regex", Value: search}, {Key: "$options", Value: "i"}}}}}})
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
	return body
}
