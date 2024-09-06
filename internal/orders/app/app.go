package app

import "github.com/Nicknamezz00/gorder/internal/orders/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	CustomerOrder query.CustomerOrder
}

