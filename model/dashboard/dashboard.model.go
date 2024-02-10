package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DashboardResponse struct {
	TotalWork            int                           `json:"totalWork" bson:"totalWork"`
	TotalWorkUnfinished  int                           `json:"totalWorkUnfinished" bson:"totalWorkUnfinished"`
	TotalEarn            int32                         `json:"totalEarn" bson:"totalEarn"`
	TotalExpense         int32                         `json:"totalExpense" bson:"totalExpense"`
	YearsReport          []GetYearReportResult         `json:"yearsReport" bson:"yearsReport"`
	CustomerWorkRatio    []CustomerWork                `json:"customerWorkRatio" bson:"customerWorkRatio"`
	CustomerProfitRatio  []CustomerMoney               `json:"customerProfitRatio" bson:"customerProfitRatio"`
	EarnAndSpendEachYear GetEarnAndSpendEachYearResult `json:"spentAndEarnEachMonth" bson:"spentAndEarnEachMonth"`
}

type GetEarnAndSpendEachYearInput struct {
	Year int `json:"year" bson:"year"`
}

type WorkResult struct {
	ID struct {
		Month int `json:"month" bson:"month"`
		Year  int `json:"year" bson:"year"`
	} `json:"_id" bson:"_id"`
	Earn float64 `json:"earn" bson:"earn"`
}
type WorkYearReportResult struct {
	ID struct {
		Month int `json:"month" bson:"month"`
		Year  int `json:"year" bson:"year"`
	} `json:"_id" bson:"_id"`
	TotalEarn float64 `json:"totalEarn" bson:"totalEarn"`
}

type ExpenseYearReportResult struct {
	ID struct {
		Month int `json:"month" bson:"month"`
		Year  int `json:"year" bson:"year"`
	} `json:"_id" bson:"_id"`
	TotalExpense float64 `json:"totalExpense" bson:"totalExpense"`
}
type ExpenseResult struct {
	ID struct {
		Month int `json:"month" bson:"month"`
		Year  int `json:"year" bson:"year"`
	} `json:"_id" bson:"_id"`
	SpendWithVat    float64 `json:"spendWithVat" bson:"spendWithVat"`
	SpendWithOutVat float64 `json:"spendWithOutVat" bson:"spendWithOutVat"`
}

type MonthEarnAndSpend struct {
	Month           int     `json:"month" bson:"month"`
	SpendWithVat    float64 `json:"spendWithVat" bson:"spendWithVat"`
	SpendWithOutVat float64 `json:"spendWithOutVat" bson:"spendWithOutVat"`
	Earn            float64 `json:"earn" bson:"earn"`
}

type GetEarnAndSpendEachYearResult struct {
	Month []MonthEarnAndSpend `json:"month" bson:"month"`
	Years []int               `json:"years" bson:"years"`
}

type GetTotalEarnInput struct {
	Year int `json:"year" bson:"year" validate:"required"`
}

type GetTotalSpendInput struct {
	Year int `json:"year" bson:"year" validate:"required"`
}

type GetTotalWorkResult struct {
	TotalWork           int `json:"totalWork" bson:"totalWork"`
	TotalWorkUnfinished int `json:"totalWorkUnfinished" bson:"totalWorkUnfinished"`
}

type GetYearReportResult struct {
	Year         int   `json:"year" bson:"year"`
	TotalEarn    int32 `json:"totalEarn" bson:"totalEarn"`
	TotalExpense int32 `json:"totalExpense" bson:"totalExpense"`
}

type GetWorkCustomerProcess struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	TotalWork int                `json:"totalWork" bson:"totalWork"`
	TotalEarn int                `json:"totalEarn" bson:"totalEarn"`
}

type GetWorkCustomerResult struct {
	CustomerWork  []CustomerWork  `json:"customerWork" bson:"customerWork"`
	CustomerMoney []CustomerMoney `json:"customerMoney" bson:"customerMoney"`
}

type CustomerWork struct {
	Name      string  `json:"name" bson:"name"`
	WorkCount int     `json:"workCount" bson:"workCount"`
	Color     string  `json:"color" bson:"color"`
	Ratio     float64 `json:"ratio" bson:"ratio"`
}

type CustomerMoney struct {
	Name      string  `json:"name" bson:"name"`
	TotalEarn int     `json:"totalEarn" bson:"totalEarn"`
	Color     string  `json:"color" bson:"color"`
	Ratio     float64 `json:"ratio" bson:"ratio"`
}

type GetTotalEarnResult struct {
	TotalEarn int32 `json:"totalEarn" bson:"totalEarn"`
}

type GetTotalSpendResult struct {
	TotalExpense int32 `json:"totalExpense" bson:"totalExpense"`
}
