package models

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID         int    `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	Weight     int    `db:"weight" json:"weight"`
	CustomerID int    `db:"customer_id" json:"customer_id"`
	WorkerID   int    `db:"worker_id" json:"worker_id"`
}

type GetOrderResponse struct {
	Name       string `db:"name" json:"name"`
	Weight     int    `db:"weight" json:"weight"`
	CustomerID int    `db:"customer_id" json:"customer_id"`
	WorkerID   int    `db:"worker_id" json:"worker_id"`
}

func GenerateRandomOrderName() string {
	words := []string{"apple", "banana", "cherry", "grape", "orange", "box", "table", "chair", "sofa", "TV", "bed"}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomIndex := r.Intn(len(words))
	randomIndex2 := r.Intn(len(words))
	randomIndex3 := r.Intn(len(words))
	return fmt.Sprintf("%s, %s, %s", words[randomIndex], words[randomIndex2], words[randomIndex3])
}
