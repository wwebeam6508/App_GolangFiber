package entity

import "time"

//struct works
type Work struct {
	ID       string    `bson:"_id"`
	Date     time.Time `bson:"date"`
	Detail   string    `bson:"detail"`
	DateEnd  time.Time `bson:"dateEnd"`
	Title    string    `bson:"title"`
	Profit   int       `bson:"profit"`
	Customer Customer  `bson:"customer"`
	Status   int       `bson:"status"`
}
