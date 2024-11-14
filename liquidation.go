package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

type (
	// Symbol is a trading symbol.
	Symbol string

	// PriceQuantity pair.
	PriceQuantity struct {
		Price         float64
		Quantity      float64
		Currency      string
		TotalUSDValue float64

		MinStep float64
		MinTick float64
	}

	// RawLiquidation is data from the table.
	RawLiquidation struct {
		OrderID   string  `json:"orderID"`
		Price     float64 `json:"price"`
		Symbol    Symbol  `json:"symbol"`
		LeavesQty float64 `json:"leavesQty"`
		Side      string  `json:"side"`
	}

	// Liquidation data (raw).
	Liquidation struct {
		PriceQuantity

		Symbol Symbol
		Side   string
	}

	// CombinedLiquidation ...
	CombinedLiquidation struct {
		Symbol Symbol
		Side   string

		Liquidations []PriceQuantity
	}
)

const (
	// MaxCombinedPositions caps the number of liquidations that can be combined into a single tweet.
	MaxCombinedPositions = 3

	// MaxUSDValueMergable ...
	MaxUSDValueMergable = 250000

	epsilon = 0.000001
)

func displayTick(value, tick float64) string {
	if tick == 0 {
		return humanize.Commaf(value)
	}

	if math.Log10(tick) < 0 {
		return humanize.CommafWithDigits(value, int(-math.Floor(math.Log10(tick))))
	}

	return humanize.Comma(int64(value))
}

func displayUSD(value float64) string {
	res := humanize.CommafWithDigits(value, 2)

	idx := strings.Index(res, ".")
	if idx == -1 {
		return res
	}

	if idx == len(res)-2 {
		return res + "0"
	}

	return res
}

// DisplayPrice using min tick.
func (pq PriceQuantity) DisplayPrice() string {
	return displayTick(pq.Price, pq.MinTick)
}

// DisplayQuantity using min step.
func (pq PriceQuantity) DisplayQuantity() string {
	return displayTick(pq.Quantity, pq.MinStep)
}

// ToCombined converts a single liquidation to a combined liquidation.
func (l Liquidation) ToCombined() CombinedLiquidation {
	return CombinedLiquidation{
		Symbol: l.Symbol,
		Side:   l.Side,
		Liquidations: []PriceQuantity{
			l.PriceQuantity,
		},
	}
}

// CanCombine returns if an addtional liquidation can be merged into an existing combined liquidation.
func (cl CombinedLiquidation) CanCombine(l Liquidation) bool {
	if cl.Side != l.Side || cl.Symbol != l.Symbol {
		return false
	}

	if len(cl.Liquidations) >= MaxCombinedPositions {
		return false
	}

	for _, l2 := range cl.Liquidations {
		if l2.TotalUSDValue > MaxUSDValueMergable {
			return false
		}
	}

	if l.TotalUSDValue > MaxUSDValueMergable {
		return false
	}

	// Acceptable merge groups
	mergeGroups := []int{1, 24999, 250000}

	for _, l2 := range cl.Liquidations {
		if sort.SearchInts(mergeGroups, int(l.Quantity)) != sort.SearchInts(mergeGroups, int(l2.Quantity)) {
			return false
		}
	}

	return true
}

// Combine an existing liquidation into the the combined liquidation.
func (cl *CombinedLiquidation) Combine(l Liquidation) error {
	if !cl.CanCombine(l) {
		return errors.New("cannot merge")
	}

	cl.Liquidations = append(cl.Liquidations, l.PriceQuantity)
	return nil
}

// String implements Stringer.
func (l Liquidation) String() string {
	var position string
	if l.Side == "Buy" {
		position = "short"
	} else {
		position = "long"
	}

	switch l.Currency {
	case "USD", "USDT":
		// Example: Liquidated short on XBTUSD: Buy 130170 @ 772.02
		return fmt.Sprintf("Liquidated %v on %v: %v %v @ %v", position, l.Symbol, l.Side, l.DisplayQuantity(), l.DisplayPrice())

	default:
		// Example: Liquidated short on ETHUSD: Buy 130170 Cont @ 772.02 (≈ $XXXX)
		if l.TotalUSDValue < epsilon {
			return fmt.Sprintf("Liquidated %v on %v: %v %v %v @ %v", position, l.Symbol, l.Side, l.DisplayQuantity(), l.Currency, l.DisplayPrice())
		} else {
			return fmt.Sprintf("Liquidated %v on %v: %v %v %v @ %v (≈ $%v)", position, l.Symbol, l.Side, l.DisplayQuantity(), l.Currency, l.DisplayPrice(), displayUSD(l.TotalUSDValue))
		}
	}

}

// String implements Stringer.
func (cl CombinedLiquidation) String() string {

	var totalValue float64
	currPrice := cl.Liquidations[0].DisplayPrice()
	samePrice := true

	for _, l := range cl.Liquidations {
		totalValue += l.TotalUSDValue
		samePrice = samePrice && l.DisplayPrice() == currPrice
	}

	var position string
	if cl.Side == "Buy" {
		position = "short"
	} else {
		position = "long"
	}

	cp := ""
	for i, l := range cl.Liquidations {
		if i > 0 {
			cp += " + "
		}
		cp += l.DisplayQuantity()
	}

	switch cl.Liquidations[0].Currency {
	case "USD", "USDT":
	default:
		cp += " " + cl.Liquidations[0].Currency
	}

	cp += " @ "
	if samePrice {
		cp += currPrice
	} else {
		for i, l := range cl.Liquidations {
			if i > 0 {
				cp += ", "
			}
			cp += l.DisplayPrice()
		}
	}

	switch cl.Liquidations[0].Currency {
	case "USD", "USDT":
		// Example: Liquidated short on XBTUSD: Buy 130170, 123450 @ 772.02, 734.01
		return fmt.Sprintf("Liquidated %v on %v: %v %s", position, cl.Symbol, cl.Side, cp)

	default:
		// Example Liquidated short on ETHUSD: Buy 100, 200 Cont @ 772.02, 734.01 (≈ $XXXX)
		return fmt.Sprintf("Liquidated %v on %v: %v %s (≈ $%v)", position, cl.Symbol, cl.Side, cp, displayUSD(totalValue))
	}
}

// CombiningDelay is the minimum time to wait for another liquidation to combine with.
// Small positions incur longer combining delays
func (l Liquidation) CombiningDelay() time.Duration {
	if l.Quantity < 1000 {
		return 30 * time.Second
	} else if l.Quantity < 25000 {
		return 20 * time.Second
	} else if l.Quantity < 125000 {
		return 15 * time.Second
	} else {
		return 10 * time.Second
	}
}

// USDValue returns the USD value of the liquidation.
func (cl CombinedLiquidation) USDValue() (total float64) {
	for _, v := range cl.Liquidations {
		total += v.TotalUSDValue
	}

	return total
}

// TotalQuantity of a combined liquidation.
func (cl CombinedLiquidation) TotalQuantity() (total float64) {
	for _, v := range cl.Liquidations {
		total += v.Quantity
	}

	return total
}

// MaxQuantity of a combined liquidation.
func (cl CombinedLiquidation) MaxQuantity() (max float64) {
	for _, v := range cl.Liquidations {
		if v.Quantity > max {
			max = v.Quantity
		}
	}

	return max
}

// MinQuantity of a combined liquidation.
func (cl CombinedLiquidation) MinQuantity() (min float64) {
	if len(cl.Liquidations) == 0 {
		return 0
	}

	min = cl.Liquidations[0].Quantity
	for _, q := range cl.Liquidations {
		if q.Quantity < min {
			min = q.Quantity
		}
	}

	return min
}
