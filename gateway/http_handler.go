package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Nicknamezz00/gorder/gateway/entry"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"github.com/Nicknamezz00/gorder/pkg/errcode"
	"github.com/Nicknamezz00/gorder/pkg/jsonutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	entry entry.OrdersEntry
}

func NewHandler(entry entry.OrdersEntry) *handler {
	return &handler{
		entry: entry,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	// embed static files
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.HandleGetOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemWithQuantity
	if err := jsonutil.ReadJSON(r, &items); err != nil {
		jsonutil.WriteWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validate(items); err != nil {
		jsonutil.WriteWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.entry.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	err2 := status.Convert(err)
	if err2 != nil {
		if err2.Code() != codes.InvalidArgument {
			jsonutil.WriteWithError(w, http.StatusBadRequest, err2.Message())
			return
		}
		jsonutil.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := &CreateOrderRequest{
		Order:         o,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s", o.CustomerID, o.ID),
	}
	jsonutil.WriteJSON(w, http.StatusOK, res)
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

func (h *handler) HandleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	orderID := r.PathValue("orderID")
	o, err := h.entry.GetOrder(r.Context(), orderID, customerID)
	err2 := status.Convert(err)
	if err2 != nil {
		if err2.Code() != codes.InvalidArgument {
			jsonutil.WriteWithError(w, http.StatusBadRequest, err2.Message())
			return
		}
		jsonutil.WriteWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jsonutil.WriteJSON(w, http.StatusOK, o)
}
