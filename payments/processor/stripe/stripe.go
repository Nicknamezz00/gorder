package stripe

import (
	"fmt"
	"log"

	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"github.com/Nicknamezz00/gorder/pkg/envutil"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

var (
	gatewayHTTPAddr = envutil.EnvString("GATEWAY_HTTP_ADDRESS", "http://127.0.0.1:8080")
)

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(o *pb.Order) (string, error) {
	log.Printf("creating payment link for order: %v", o)
	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, o.CustomerID, o.ID)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html?customerID=%s&orderID=%s", gatewayHTTPAddr, o.CustomerID, o.ID)
	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, it := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(it.PriceID),
			Quantity: stripe.Int64(int64(it.Quantity)),
		})
	}
	meta := map[string]string{
		"customerID": o.CustomerID,
		"orderID":    o.ID,
	}
	params := &stripe.CheckoutSessionParams{
		Metadata:   meta,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, err
}
