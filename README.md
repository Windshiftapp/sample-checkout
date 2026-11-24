# sample-checkout

> **This is a fictional sample repository.** It exists to give the public
> Windshift demo a small, inspectable bug-fix story: a customer report, a real
> pull request, and a regression test. It is **not** Windshift core and shares
> no code with it. The shop, the products, and the customer are invented.

A tiny checkout service for "Acme Storefront". It prices a cart, applies an
optional percentage coupon, and charges the customer through a payment gateway.

## Layout

| Path | What |
|---|---|
| `checkout/cart.go` | Cart line items, subtotals, and coupon math |
| `checkout/gateway.go` | Payment gateway interface |
| `checkout/service.go` | `Checkout`: applies the coupon, then charges with retries |
| `checkout/service_test.go` | Tests |

## Run the tests

```bash
go test ./...
```

## The story

A customer reported that an order is charged the **full price** when the
payment is retried after a transient gateway failure. The coupon is applied to
the first attempt but lost on the retry. The fix and its regression test live in
the merged pull request for this repository.

## License

Sample code, provided as-is for demonstration purposes.
