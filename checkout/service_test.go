package checkout

import (
	"errors"
	"testing"
)

// recordingGateway records every amount it is asked to charge and fails the
// first `failures` attempts with a transient error.
type recordingGateway struct {
	failures int
	charges  []int64
}

func (g *recordingGateway) Charge(amountCents int64, _ string) (string, error) {
	g.charges = append(g.charges, amountCents)
	if len(g.charges) <= g.failures {
		return "", ErrTransient
	}
	return "ch_test", nil
}

func TestCheckoutAppliesCoupon(t *testing.T) {
	gateway := &recordingGateway{}
	service := &Service{Gateway: gateway}
	cart := Cart{Items: []LineItem{{SKU: "coffee", UnitPriceCents: 2000, Quantity: 1}}}
	coupon := &Coupon{Code: "SAVE20", PercentOff: 20}

	receipt, err := service.Checkout(cart, coupon)
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}
	if want := int64(1600); receipt.ChargedCents != want {
		t.Fatalf("charged %d cents, want %d", receipt.ChargedCents, want)
	}
	if len(gateway.charges) != 1 || gateway.charges[0] != 1600 {
		t.Fatalf("gateway charges = %v, want [1600]", gateway.charges)
	}
}

func TestCheckoutRetriesTransientFailure(t *testing.T) {
	gateway := &recordingGateway{failures: 1}
	service := &Service{Gateway: gateway}
	cart := Cart{Items: []LineItem{{SKU: "mug", UnitPriceCents: 1500, Quantity: 2}}}

	receipt, err := service.Checkout(cart, nil)
	if err != nil {
		t.Fatalf("Checkout: %v", err)
	}
	if want := int64(3000); receipt.ChargedCents != want {
		t.Fatalf("charged %d cents, want %d", receipt.ChargedCents, want)
	}
	if receipt.Attempts != 2 {
		t.Fatalf("attempts = %d, want 2", receipt.Attempts)
	}
}

func TestCheckoutStopsOnPermanentFailure(t *testing.T) {
	service := &Service{Gateway: permanentFailureGateway{}}
	cart := Cart{Items: []LineItem{{SKU: "mug", UnitPriceCents: 1500, Quantity: 1}}}

	if _, err := service.Checkout(cart, nil); !errors.Is(err, errDeclined) {
		t.Fatalf("Checkout error = %v, want %v", err, errDeclined)
	}
}

var errDeclined = errors.New("card declined")

type permanentFailureGateway struct{}

func (permanentFailureGateway) Charge(int64, string) (string, error) {
	return "", errDeclined
}
