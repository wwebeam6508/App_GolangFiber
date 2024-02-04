package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

type GetCustomerInput struct {
	Page         int    `json:"page" bson:"page"`
	PageSize     int    `json:"pageSize" bson:"pageSize"`
	SortTitle    string `json:"sortTitle" bson:"sortTitle"`
	SortType     string `json:"sortType" bson:"sortType"`
	Search       string `json:"search" bson:"search"`
	SearchFilter string `json:"searchFilter" bson:"searchFilter"`
}

type SearchPipeline struct {
	Search         string `json:"search" bson:"search"`
	SearchPipeline bson.A `json:"searchPipeline" bson:"searchPipeline"`
}

type GetCustomerResult struct {
	CustomerID string `json:"customerID" bson:"customerID"`
	Name       string `json:"name" bson:"name"`
	Address    string `json:"address" bson:"address"`
	TaxID      string `json:"taxID" bson:"taxID"`
}

type GetCustomerByIDInput struct {
	CustomerID string `json:"customerID" bson:"customerID" validate:"required"`
}

type GetCustomerByIDResult struct {
	CustomerID string   `json:"customerID" bson:"customerID"`
	Name       string   `json:"name" bson:"name"`
	Address    string   `json:"address" bson:"address"`
	TaxID      string   `json:"taxID" bson:"taxID"`
	Emails     []string `json:"emails" bson:"emails"`
	Phones     []string `json:"phones" bson:"phones"`
}

type AddCustomerInput struct {
	Name    string   `json:"name" bson:"name" validate:"required"`
	Address string   `json:"address" bson:"address"`
	TaxID   string   `json:"taxID" bson:"taxID"`
	Emails  []string `json:"emails" bson:"emails"`
	Phones  []string `json:"phones" bson:"phones"`
}

type UpdateCustomerID struct {
	CustomerID string `json:"customerID" bson:"customerID" validate:"required"`
}

type UpdateCustomerInput struct {
	Name    string   `json:"name" bson:"name"`
	Address string   `json:"address" bson:"address"`
	TaxID   string   `json:"taxID" bson:"taxID"`
	Emails  []string `json:"emails" bson:"emails" validate:"email"`
	Phones  []string `json:"phones" bson:"phones" `
}
