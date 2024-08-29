package entry

import "context"

type PaymentEntry interface {
	UpdateOrderAfterPaid(ctx context.Context, orderID, paymentLink string) error
}
