package checkout

import "errors"

// ErrTransient marks a gateway failure that is safe to retry.
var ErrTransient = errors.New("temporary gateway failure")

// Gateway charges a customer. Implementations talk to the payment provider.
type Gateway interface {
	// Charge attempts to charge amountCents. idempotencyKey is stable across
	// retries of the same checkout so the provider does not double-charge.
	Charge(amountCents int64, idempotencyKey string) (chargeID string, err error)
}
