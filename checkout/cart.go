package checkout

// LineItem is one line of a shopping cart.
type LineItem struct {
	SKU            string
	UnitPriceCents int64
	Quantity       int
}

// Cart holds the items a customer is buying.
type Cart struct {
	Items []LineItem
}

// SubtotalCents is the cart total before any coupon is applied.
func (c Cart) SubtotalCents() int64 {
	var total int64
	for _, item := range c.Items {
		total += item.UnitPriceCents * int64(item.Quantity)
	}
	return total
}

// Coupon is a percentage discount code.
type Coupon struct {
	Code       string
	PercentOff int
}

// DiscountCents is the amount taken off a given subtotal.
func (c Coupon) DiscountCents(subtotalCents int64) int64 {
	if c.PercentOff <= 0 {
		return 0
	}
	if c.PercentOff >= 100 {
		return subtotalCents
	}
	return subtotalCents * int64(c.PercentOff) / 100
}
