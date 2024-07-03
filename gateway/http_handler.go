package main

import (
	"errors"
	"net/http"

	common "github.com/Nicknamezz00/gorder-common"
	pb "github.com/Nicknamezz00/gorder-common/api"
	"github.com/Nicknamezz00/gorder-common/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{
		client: client,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate(items); err != nil {
		common.WriteWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	err2 := status.Convert(err)
	if err2 != nil {
		if err2.Code() != codes.InvalidArgument {
			common.WriteWithError(w, http.StatusBadRequest, err2.Message())
			return
		}
		common.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, o)
}

func validate(items []*pb.ItemWithQuantity) error {
	if len(items) == 0 {
		return errcode.ErrNoItems
	}
	for _, i := range items {
		if i.ID == "" {
			return errors.New("item id is required")
		}
		if i.Quantity <= 0 {
			return errors.New("item quantity must be valid")
		}
	}
	return nil
}
