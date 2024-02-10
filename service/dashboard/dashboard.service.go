package service

import (
	"PBD_backend_go/common"
	"PBD_backend_go/configuration"
	model "PBD_backend_go/model/dashboard"
	"context"
	"math"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetSpentAndEarnEachYear(year *int) (*model.GetEarnAndSpendEachYearResult, error) {
	years, err := GetYearsList()
	if err != nil {
		return nil, err
	}
	if common.IsEmpty(year) {
		year = &years[0]
	}
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	workResult, err := getWorkEarn(*year, workRef)
	if err != nil {
		return nil, err
	}
	expenseRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	expenseResult, err := getExpenseSpend(*year, expenseRef)
	if err != nil {
		return nil, err
	}
	var months []model.MonthEarnAndSpend
	for i := 0; i < 12; i++ {
		months = append(months, model.MonthEarnAndSpend{
			Month:           i + 1,
			SpendWithVat:    0,
			SpendWithOutVat: 0,
			Earn:            0,
		})
	}

	for _, work := range *workResult {
		months[work.ID.Month-1].Earn = work.Earn
	}
	for _, expense := range *expenseResult {
		months[expense.ID.Month-1].SpendWithVat = expense.SpendWithVat
		months[expense.ID.Month-1].SpendWithOutVat = expense.SpendWithOutVat
	}

	return &model.GetEarnAndSpendEachYearResult{
		Month: months,
		Years: years,
	}, nil

}

func GetTotalEarn(year *int) (*int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	pipelineWork := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "dateEnd", Value: bson.D{
				{Key: "$ne", Value: nil},
			}},
			{Key: "status", Value: 1},
		}}},
	}
	if !common.IsEmpty(year) {
		pipelineWork = append(pipelineWork, bson.D{{Key: "$match", Value: bson.D{
			{Key: "dateEnd", Value: bson.D{
				{Key: "$ne", Value: nil},
			}},
			{Key: "status", Value: 1},
		}}})
	}
	if !common.IsEmpty(year) {
		start := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(*year, 12, 31, 23, 59, 59, 0, time.UTC)
		pipelineWork = append(pipelineWork, bson.D{{Key: "$match", Value: bson.D{{Key: "dateEnd", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}}}}})
	}

	pipelineWork = append(pipelineWork, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: nil},
		{Key: "totalEarn", Value: bson.D{{Key: "$sum", Value: "$profit"}}},
	}}})
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	var totalEarn bson.M
	if workCursor.Next(context.Background()) {
		workCursor.Decode(&totalEarn)
	}
	earn := totalEarn["totalEarn"].(int32)
	return &earn, nil
}

func GetTotalExpense(year *int) (*int32, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	expenseRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")
	pipelineExpense := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
	}
	if !common.IsEmpty(year) {
		start := time.Date(*year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(*year, 12, 31, 23, 59, 59, 0, time.UTC)
		pipelineExpense = append(pipelineExpense, bson.D{{Key: "$match", Value: bson.D{{Key: "date", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}}}}})
	}
	totalExpenseReduce := bson.D{{Key: "$reduce", Value: bson.D{
		{Key: "input", Value: "$lists"},
		{Key: "initialValue", Value: 0},
		{Key: "in", Value: bson.D{
			{Key: "$add", Value: bson.A{"$$value", "$$this.price"}},
		},
		},
	}}}
	pipelineExpense = append(pipelineExpense, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: nil},
		{Key: "totalExpense", Value: bson.D{{Key: "$sum", Value: totalExpenseReduce}}},
	}}})

	expenseCursor, err := expenseRef.Aggregate(context.Background(), pipelineExpense)
	if err != nil {
		return nil, err
	}
	var totalExpense bson.M
	if expenseCursor.Next(context.Background()) {
		expenseCursor.Decode(&totalExpense)
	}
	expense := int32(totalExpense["totalExpense"].(float64))
	return &expense, nil
}

func GetYearsReport() ([]model.GetYearReportResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	expenseRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")

	pipelineWork := bson.A{
		//status eq 1
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: 1},
			{Key: "dateEnd", Value: bson.D{{Key: "$ne", Value: nil}}},
		}}},
		// dateEnd not null
		// group _id :{ year: { $year: "$dateEnd" } }
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "year", Value: bson.D{{Key: "$year", Value: "$dateEnd"}}}}},
			{Key: "totalEarn", Value: bson.D{{Key: "$sum", Value: "$profit"}}},
		}}},
	}
	pipelineExpense := bson.A{
		//status eq 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
		// group _id :{ year: { $year: "$date" } }
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "year", Value: bson.D{{Key: "$year", Value: "$date"}}}}},
			{Key: "totalExpense", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$reduce", Value: bson.D{
					{Key: "input", Value: "$lists"},
					{Key: "initialValue", Value: 0},
					{Key: "in", Value: bson.D{
						{Key: "$add", Value: bson.A{"$$value", "$$this.price"}},
					},
					},
				}},
			}}},
			}}}},
	}
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	expenseCursor, err := expenseRef.Aggregate(context.Background(), pipelineExpense)
	if err != nil {
		return nil, err
	}
	var workResult []model.WorkYearReportResult
	var expenseResult []model.ExpenseYearReportResult
	if err = workCursor.All(context.Background(), &workResult); err != nil {
		return nil, err
	}
	if err = expenseCursor.All(context.Background(), &expenseResult); err != nil {
		return nil, err
	}
	var yearsReport []model.GetYearReportResult
	for _, work := range workResult {
		year := work.ID.Year
		yearsReport = append(yearsReport, model.GetYearReportResult{
			Year:         year,
			TotalEarn:    int32(work.TotalEarn),
			TotalExpense: 0,
		})
	}
	for _, expense := range expenseResult {
		year := expense.ID.Year
		index := common.FindIndex(yearsReport, func(y model.GetYearReportResult) bool {
			return y.Year == year
		})
		if index != -1 {
			yearsReport[index].TotalExpense = int32(expense.TotalExpense)
		} else {
			yearsReport = append(yearsReport, model.GetYearReportResult{
				Year:         year,
				TotalEarn:    0,
				TotalExpense: int32(expense.TotalExpense),
			})
		}
	}
	return yearsReport, nil
}

func GetTotalWork() (*model.GetTotalWorkResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	pipelineWork := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
		//group totalWork : { $sum : 1 }
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalWork", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "totalWorkUnfinished", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$cond", Value: bson.A{
					bson.D{{Key: "$eq", Value: bson.A{"$dateEnd", nil}}},
					1,
					0,
				},
				},
			}}},
			},
		}}},
	}
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	var totalWork model.GetTotalWorkResult
	if workCursor.Next(context.Background()) {
		workCursor.Decode(&totalWork)
	}
	return &totalWork, nil
}

func GetWorkCustomer() (*model.GetWorkCustomerResult, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	pipelineWork := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: 1},
			{Key: "dateEnd", Value: bson.D{{Key: "$ne", Value: nil}}},
		}}},
		//lookup customer
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "customers"},
			{Key: "localField", Value: "customer.$id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "customer"},
		}}},
		//unwind customer
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$customer"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		//group by customer
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$customer._id"},
			{Key: "name", Value: bson.D{{Key: "$first", Value: "$customer.name"}}},
			{Key: "totalWork", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "totalEarn", Value: bson.D{{Key: "$sum", Value: "$profit"}}},
		}}}}
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	var workCustomerResult []model.GetWorkCustomerProcess
	if err = workCursor.All(context.Background(), &workCustomerResult); err != nil {
		return nil, err
	}
	var customerWork []model.CustomerWork
	var customerMoney []model.CustomerMoney
	var totalWorks int
	var totalEarns int
	for _, work := range workCustomerResult {
		color := common.GenerateRandomColor()
		customerWork = append(customerWork, model.CustomerWork{
			Name:      work.Name,
			WorkCount: work.TotalWork,
			Color:     color,
			Ratio:     0,
		})
		totalWorks += work.TotalWork
		customerMoney = append(customerMoney, model.CustomerMoney{
			Name:      work.Name,
			TotalEarn: work.TotalEarn,
			Color:     color,
			Ratio:     0,
		})
		totalEarns += work.TotalEarn
	}

	//calculate ratio and fixed 2 decimal
	for i, work := range customerWork {
		customerWork[i].Ratio = math.Round(float64(work.WorkCount) / float64(totalWorks) * 100)
	}
	for i, money := range customerMoney {
		customerMoney[i].Ratio = math.Round(float64(money.TotalEarn) / float64(totalEarns) * 100)
	}
	return &model.GetWorkCustomerResult{

		CustomerWork:  customerWork,
		CustomerMoney: customerMoney,
	}, nil
}

func getWorkEarn(year int, workRef *mongo.Collection) (*[]model.WorkResult, error) {
	pipelineWork := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
		//match where has date end
		bson.D{{Key: "$match", Value: bson.D{{Key: "dateEnd", Value: bson.D{{Key: "$ne", Value: nil}}}}}},
	}
	if !common.IsEmpty(year) {
		start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
		pipelineWork = append(pipelineWork, bson.D{{Key: "$match", Value: bson.D{{Key: "dateEnd", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}}}}})
	}
	pipelineWork = append(pipelineWork, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "month", Value: bson.D{{Key: "$month", Value: "$dateEnd"}}}}},
		{Key: "earn", Value: bson.D{{Key: "$sum", Value: "$profit"}}},
	}}})
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	var workResult []model.WorkResult
	if err = workCursor.All(context.Background(), &workResult); err != nil {
		return nil, err
	}
	return &workResult, nil
}

func getExpenseSpend(year int, expenseRef *mongo.Collection) (*[]model.ExpenseResult, error) {
	pipelineExpense := bson.A{
		//match status = 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
	}
	if !common.IsEmpty(year) {
		start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
		pipelineExpense = append(pipelineExpense, bson.D{{Key: "$match", Value: bson.D{{Key: "date", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}}}}})
	}

	spendWithVat := bson.D{{Key: "$reduce", Value: bson.D{
		{Key: "input", Value: "$lists"},
		{Key: "initialValue", Value: 0},
		{Key: "in", Value: bson.D{
			{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$gt", Value: bson.A{"$currentVat", 0}}},
				bson.D{{Key: "$add", Value: bson.A{"$$value", "$$this.price"}}},
				"$$value",
			}},
		},
		},
	}}}
	spendWithOutVat := bson.D{{Key: "$reduce", Value: bson.D{
		{Key: "input", Value: "$lists"},
		{Key: "initialValue", Value: 0},
		{Key: "in", Value: bson.D{
			{Key: "$cond", Value: bson.A{
				bson.D{{Key: "$eq", Value: bson.A{"$currentVat", 0}}},
				bson.D{{Key: "$add", Value: bson.A{"$$value", "$$this.price"}}},
				"$$value",
			}},
		},
		},
	}}}
	pipelineExpense = append(pipelineExpense, bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "month", Value: bson.D{{Key: "$month", Value: "$date"}}}}},
		{Key: "spendWithVat", Value: bson.D{{Key: "$sum", Value: spendWithVat}}},
		{Key: "spendWithOutVat", Value: bson.D{{Key: "$sum", Value: spendWithOutVat}}},
	}}})

	expenseCursor, err := expenseRef.Aggregate(context.Background(), pipelineExpense)
	if err != nil {
		return nil, err
	}
	var expenseResult []model.ExpenseResult
	if err = expenseCursor.All(context.Background(), &expenseResult); err != nil {
		return nil, err
	}
	return &expenseResult, nil
}

func GetYearsList() ([]int, error) {
	coll, err := configuration.ConnectToMongoDB()
	defer coll.Disconnect(context.Background())
	if err != nil {
		return nil, err
	}
	workRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("works")
	expenseRef := coll.Database(os.Getenv("MONGO_DB_NAME")).Collection("expenses")

	pipelineWork := bson.A{
		//status eq 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
		// group _id :{ year: { $year: "$date" } }
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "year", Value: bson.D{{Key: "$year", Value: "$date"}}}}},
		}}},
	}
	pipelineExpense := bson.A{
		//status eq 1
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: 1}}}},
		// group _id :{ year: { $year: "$date" } }
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "year", Value: bson.D{{Key: "$year", Value: "$date"}}}}},
		}}},
	}
	workCursor, err := workRef.Aggregate(context.Background(), pipelineWork)
	if err != nil {
		return nil, err
	}
	expenseCursor, err := expenseRef.Aggregate(context.Background(), pipelineExpense)
	if err != nil {
		return nil, err
	}
	var workYears []model.WorkResult
	var expenseYears []model.ExpenseResult
	if err = workCursor.All(context.Background(), &workYears); err != nil {
		return nil, err
	}
	if err = expenseCursor.All(context.Background(), &expenseYears); err != nil {
		return nil, err
	}
	var years []int
	for _, work := range workYears {
		years = append(years, work.ID.Year)
	}
	for _, expense := range expenseYears {
		if !common.Contains(years, expense.ID.Year) {
			years = append(years, expense.ID.Year)
		}
	}
	//sort years desc
	common.SortIntDesc(&years)
	return years, nil
}
