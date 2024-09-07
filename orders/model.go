package main

import (
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID  string             `bson:"customerID,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"paymentLink,omitempty"`
	Items       []*pb.Item         `bson:"items,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		ID:          o.ID.Hex(),
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
