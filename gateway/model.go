package main

import pb "github.com/Nicknamezz00/gorder/pkg/api"

type CreateOrderRequest struct {
	Order         *pb.Order `"json": order`
	RedirectToURL string    `"json": redirectToURL`
}
