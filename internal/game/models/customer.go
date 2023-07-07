package models

import "github.com/go-playground/validator/v10"

var (
	validate *validator.Validate
)

const (
	CustomerRole = "customer"
)

func init() {
	validate = validator.New()
}

type Customer struct {
	UserID       int `db:"user_id" json:"user_id"`
	StartCapital int `db:"start_capital" json:"start_capital"`
}

type StartGameRequest struct {
	WorkerIDs []int `json:"workerIDs"`
	OrderID   int   `json:"orderID"`
}

type GetCustomerInfoResponse struct {
	CustomerStartCapital int      `json:"start_capital"`
	Workers              []Worker `json:"workers"`
}

type GetAvaliableOrdersForCustomerResponse struct {
	ID     int    `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	Weight int    `db:"weight" json:"weight"`
}
