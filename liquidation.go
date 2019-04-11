package main

import (
	"errors"
	"fmt"
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
	// MaxUSDMergable liquidations larger than this size will not be merged.
	MaxUSDMergable = 750000

	// MaxCombinedPositions caps the number of liquidations that can be combined into a single tweet.
	MaxCombinedPositions = 5
)

// IsUSDContract returns if the contract is in USD.
func (l Liquidation) IsUSDContract() bool {
	return strings.HasPrefix(string(l.Symbol), "XBT")
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

	if l.IsUSDContract() && l.Quantity > MaxUSDMergable {
		return false
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

// USDValue returns the USD value of the liquidation.
func (l Liquidation) USDValue() int64 {
	// All XBT contracts are in USD
	if l.IsUSDContract() {
		return l.Quantity
	}

	return 0
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

// IsUSDContract return if the contract is in USD sizes.
func (cl CombinedLiquidation) IsUSDContract() bool {
	return strings.HasPrefix(string(cl.Symbol), "XBT")
}

// USDValue returns the USD value of the liquidation.
func (cl CombinedLiquidation) USDValue() int64 {
	// All XBT contracts are in USD
	if cl.IsUSDContract() {
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
