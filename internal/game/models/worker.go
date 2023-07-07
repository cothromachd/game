package models

const (
	WorkerRole = "worker"
)

type Worker struct {
	UserID    int     `db:"user_id" json:"user_id"`
	MaxWeight int     `db:"max_weight" json:"max_weight"`
	IsDrunk   bool    `db:"is_drunk" json:"is_drunk"`
	Fatigue   float64 `db:"fatigue" json:"fatigue"`
	Salary    int     `db:"salary" json:"salary"`
}

type GetAndCreateWorkerInfo struct {
	MaxWeight int     `db:"max_weight" json:"max_weight"`
	IsDrunk   bool    `db:"is_drunk" json:"is_drunk"`
	Fatigue   float64 `db:"fatigue" json:"fatigue"`
	Salary    int     `db:"salary" json:"salary"`
}

type GetCompletedWorkerOrdersResponse struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Weight   int    `db:"weight" json:"weight"`
	WorkerID int    `db:"worker_id" json:"worker_id"`
}
