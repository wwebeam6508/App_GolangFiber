package entity

//struct customers

type Customer struct {
	ID      string   `bson:"_id"`
	Emails  []string `bson:"emails"`
	Address string   `bson:"address"`
	TaxID   string   `bson:"taxID"`
	Name    string   `bson:"name"`
	Phones  []string `bson:"phones"`
	Status  int      `bson:"status"`
}
