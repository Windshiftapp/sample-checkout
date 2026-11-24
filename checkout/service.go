package checkout

import (
	"errors"
	"fmt"
)

const maxAttempts = 3

// Receipt describes a successful charge.
type Receipt struct {
	ChargeID     string
	ChargedCents int64
	Attempts     int
}

// Service prices a cart and charges the customer.
type Service struct {
	Gateway Gateway
}

// Checkout applies the coupon (when present) and charges the customer. A
// transient gateway failure is retried up to maxAttempts times.
func (s *Service) Checkout(cart Cart, coupon *Coupon) (Receipt, error) {
	amount := cart.SubtotalCents()
	if coupon != nil {
		amount -= coupon.DiscountCents(cart.SubtotalCents())
	}
	return s.chargeWithRetry(amount)
}

func (s *Service) chargeWithRetry(amount int64) (Receipt, error) {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		chargeID, err := s.Gateway.Charge(amount, "checkout-1")
		if err == nil {
			return Receipt{ChargeID: chargeID, ChargedCents: amount, Attempts: attempt}, nil
		}
		if !errors.Is(err, ErrTransient) {
			return Receipt{}, err
		}
		lastErr = err
		// Retry the same quoted amount. The coupon was already applied when the
		// amount was first computed; recomputing it here would drop the discount.
	}
	return Receipt{}, fmt.Errorf("payment failed after %d attempts: %w", maxAttempts, lastErr)
}
