package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

type (
	// Symbol is a trading symbol.
	Symbol string

	// PriceQuantity pair.
	PriceQuantity struct {
		Price    float64
		Quantity int64
	}

	// Liquidation data (raw).
	Liquidation struct {
		PriceQuantity
		Symbol Symbol
		Side   string

		AmendUp bool
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
)

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
		if l2.Quantity > cl.Symbol.MaxQuantityMergable() {
			return false
		}
	}

	if l.Quantity > l.Symbol.MaxQuantityMergable() {
		return false
	}

	// Acceptable merge groups
	mergeGroups := []int{1, 2499, 24999, 250000}

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

	// Liquidated short on XBTUSD: Buy 130170 @ 772.02
	return fmt.Sprintf("Liquidated %v on %v: %v %v @ %v", position, l.Symbol, l.Side, humanize.Comma(l.Quantity), l.Price)
}

// String implements Stringer.
func (cl CombinedLiquidation) String() string {
	var position string
	if cl.Side == "Buy" {
		position = "short"
	} else {
		position = "long"
	}

	cp := ""
	for i, pc := range cl.Liquidations {
		if i > 0 {
			cp += " + "
		}
		cp += humanize.Comma(pc.Quantity)
	}

	cp += " @ "
	for i, pc := range cl.Liquidations {
		if i > 0 {
			cp += ", "
		}
		cp += fmt.Sprint(pc.Price)
	}

	// Liquidated short on XBTUSD: Buy 130170, 123450 @ 772.02, 734.01
	return fmt.Sprintf("Liquidated %v on %v: %v %s", position, cl.Symbol, cl.Side, cp)
}

// IsUSDContract returns the USD value of this contract.
func (s Symbol) IsUSDContract() bool {
	return strings.HasPrefix(string(s), "XBT")
}

// MaxQuantityMergable returns the maximum size mergable for this symbol.
func (s Symbol) MaxQuantityMergable() int64 {
	switch {
	case s.IsUSDContract():
		return 250000
	case s == "ETHUSD":
		return 500000
	case strings.HasPrefix(string(s), "ADA"):
		return 5000000
	case strings.HasPrefix(string(s), "TRX"):
		return 5000000
	default:
		return 1000000
	}
}

// USDValue returns the USD value of the liquidation.
func (cl CombinedLiquidation) USDValue() int64 {
	if cl.Symbol.IsUSDContract() {
		total := int64(0)
		for _, x := range cl.Liquidations {
			total += x.Quantity
		}

		return total
	}

	return 0
}

// TotalQuantity of a combined liquidation.
func (cl CombinedLiquidation) TotalQuantity() (total int64) {
	for _, q := range cl.Liquidations {
		total += q.Quantity
	}

	return total
}

// MaxQuantity of a combined liquidation.
func (cl CombinedLiquidation) MaxQuantity() (max int64) {
	for _, q := range cl.Liquidations {
		if q.Quantity > max {
			max = q.Quantity
		}
	}

	return max
}
