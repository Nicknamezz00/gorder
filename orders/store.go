package main

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "orders"
	CollName = "orders"
)

type OrderStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(ctx context.Context, orderID, customerID string) (*Order, error)
	Update(ctx context.Context, id string, o *pb.Order) error
}

type Store struct {
	db *mongo.Client
}

// deprecated
var inMemoryStore = make([]*pb.Order, 0)

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, o Order) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)
	newOrder, err := col.InsertOne(ctx, o)
	id := newOrder.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *Store) Get(ctx context.Context, orderID, customerID string) (*Order, error) {
	//for _, o := range inMemoryStore {
	//	if o.ID == orderID && o.CustomerID == customerID {
	//		return o, nil
	//	}
	//}
	//return nil, errcode.ErrOrderNotFound
	col := s.db.Database(DbName).Collection(CollName)
	oID, _ := primitive.ObjectIDFromHex(orderID)
	var o Order
	err := col.FindOne(ctx, bson.M{
		"_id":        oID,
		"customerID": customerID,
	}).Decode(&o)
	return &o, err
}

func (s *Store) Update(ctx context.Context, orderID string, o *pb.Order) error {
	//for i, v := range inMemoryStore {
	//	if v.ID == orderID {
	//		inMemoryStore[i].Status = o.Status
	//		inMemoryStore[i].PaymentLink = o.PaymentLink
	//		log.Printf("Order %s Updated! new status: %s", orderID, o.Status)
	//		return nil
	//	}
	//}
	//return errcode.ErrOrderNotFound
	col := s.db.Database(DbName).Collection(CollName)
	oID, _ := primitive.ObjectIDFromHex(orderID)
	_, err := col.UpdateOne(ctx,
		bson.M{"_id": oID},
		bson.M{"$set": bson.M{
			"paymentLink": o.PaymentLink,
			"status":      o.Status,
		}})
	return err
}
