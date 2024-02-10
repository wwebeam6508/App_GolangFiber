package controller

import (
	"PBD_backend_go/commonentity"
	"PBD_backend_go/exception"
	model "PBD_backend_go/model/dashboard"
	service "PBD_backend_go/service/dashboard"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardController(c *fiber.Ctx) error {
	earnAndSpendEachYearChan, customerRatioChan, totalEarnChan, totalExpenseChan, wholeTotalWorkChan, yearsReportChan, errChan := make(chan model.GetEarnAndSpendEachYearResult, 1), make(chan model.GetWorkCustomerResult, 1), make(chan int32, 1), make(chan int32, 1), make(chan model.GetTotalWorkResult, 1), make(chan []model.GetYearReportResult, 1), make(chan error, 7)
	go func() {
		earnAndSpendEachYear, err := service.GetSpentAndEarnEachYear(nil)
		if err != nil {
			errChan <- err
			earnAndSpendEachYearChan <- model.GetEarnAndSpendEachYearResult{}
			return
		}
		earnAndSpendEachYearChan <- *earnAndSpendEachYear
		errChan <- nil
	}()
	go func() {
		customerRatio, err := service.GetWorkCustomer()
		if err != nil {
			errChan <- err
			customerRatioChan <- model.GetWorkCustomerResult{}
			return
		}
		customerRatioChan <- *customerRatio
		errChan <- nil
	}()
	go func() {
		totalEarn, err := service.GetTotalEarn(nil)
		if err != nil {
			errChan <- err
			totalEarnChan <- 0
			return
		}
		totalEarnChan <- *totalEarn
		errChan <- nil
	}()
	go func() {
		totalExpense, err := service.GetTotalExpense(nil)
		if err != nil {
			errChan <- err
			totalExpenseChan <- 0
			return
		}
		totalExpenseChan <- *totalExpense
		errChan <- nil
	}()
	go func() {
		wholeTotalWork, err := service.GetTotalWork()
		if err != nil {
			errChan <- err
			wholeTotalWorkChan <- model.GetTotalWorkResult{}
			return
		}
		wholeTotalWorkChan <- *wholeTotalWork
		errChan <- nil
	}()
	go func() {
		yearsReport, err := service.GetYearsReport()
		if err != nil {
			errChan <- err
			yearsReportChan <- []model.GetYearReportResult{}
			return
		}
		yearsReportChan <- yearsReport
		errChan <- nil
	}()
	err := <-errChan
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	earnAndSpendEachYear := <-earnAndSpendEachYearChan
	customerRatio := <-customerRatioChan
	totalEarn := <-totalEarnChan
	totalExpense := <-totalExpenseChan
	wholeTotalWork := <-wholeTotalWorkChan
	yearsReport := <-yearsReportChan

	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data: model.DashboardResponse{
			TotalWork:            wholeTotalWork.TotalWork,
			TotalWorkUnfinished:  wholeTotalWork.TotalWorkUnfinished,
			TotalEarn:            totalEarn,
			TotalExpense:         totalExpense,
			YearsReport:          yearsReport,
			CustomerWorkRatio:    customerRatio.CustomerWork,
			CustomerProfitRatio:  customerRatio.CustomerMoney,
			EarnAndSpendEachYear: earnAndSpendEachYear,
		},
	})
}

func GetEarnAndSpendEachYearController(c *fiber.Ctx) error {
	var query model.GetEarnAndSpendEachYearInput
	if err := c.QueryParser(&query); err != nil {
		return exception.ErrorHandler(c, err)
	}
	earnAndSpendEachYear, err := service.GetSpentAndEarnEachYear(&query.Year)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    earnAndSpendEachYear,
	})
}

func GetTotalEarnController(c *fiber.Ctx) error {
	totalEarn, err := service.GetTotalEarn(nil)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    model.GetTotalEarnResult{TotalEarn: *totalEarn},
	})
}

func GetTotalExpenseController(c *fiber.Ctx) error {
	totalExpense, err := service.GetTotalExpense(nil)
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    model.GetTotalSpendResult{TotalExpense: *totalExpense},
	})
}

func GetTotalWorkController(c *fiber.Ctx) error {
	wholeTotalWork, err := service.GetTotalWork()
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    wholeTotalWork,
	})
}

func GetYearsReportController(c *fiber.Ctx) error {
	yearsReport, err := service.GetYearsReport()
	if err != nil {
		return exception.ErrorHandler(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(commonentity.GeneralResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    yearsReport,
	})
}


