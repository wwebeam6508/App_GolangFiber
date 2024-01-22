package entity

import (
	"time"
)

type Ref struct {
	Ref string `bson:"$ref"`
	ID  string `bson:"$id"`
}

type List struct {
	ID    string `bson:"_id"`
	Price int    `bson:"price"`
	Title string `bson:"title"`
}

type Expense struct {
	ID          string    `bson:"_id"`
	Date        time.Time `bson:"date"`
	WorkRef     Ref       `bson:"workRef"`
	Lists       []List    `bson:"lists"`
	CurrentVat  int       `bson:"currentVat"`
	Detail      string    `bson:"detail"`
	Title       string    `bson:"title"`
	Status      int       `bson:"status"`
	CustomerRef Ref       `bson:"customerRef"`
}
