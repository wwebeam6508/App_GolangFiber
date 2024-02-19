package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ref struct {
	Ref string             `bson:"$ref" json:"$ref"`
	ID  primitive.ObjectID `bson:"$id" json:"$id"`
}

type List struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Price float64            `bson:"price" json:"price"`
	Title string             `bson:"title" json:"title"`
}

type Expense struct {
	ID          string    `bson:"id" json:"id"`
	Date        time.Time `bson:"date" json:"date"`
	WorkRef     Ref       `bson:"workRef" json:"workRef"`
	Lists       []List    `bson:"lists" json:"lists"`
	CurrentVat  int       `bson:"currentVat" json:"currentVat"`
	Detail      string    `bson:"detail" json:"detail"`
	Title       string    `bson:"title" json:"title"`
	Status      int       `bson:"status" json:"status"`
	CustomerRef Ref       `bson:"customerRef" json:"customerRef"`
}
