package main

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

type (
	// Symbol is a trading symbol.
	Symbol string

	// Liquidation data.
	Liquidation struct {
		Price    float64
		Quantity int64
		Symbol   Symbol
		Side     string
	}
)

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
	if strings.HasPrefix(string(l.Symbol), "XBT") {
		return l.Quantity
	}

	// Contract value is 100.00, so it is about right
	if strings.HasPrefix(string(l.Symbol), "XBJ") {
		return l.Quantity
	}

	// Contract value is 10.00, it is not quite right (it's about 7) but close enough
	if strings.HasPrefix(string(l.Symbol), "XBC") {
		return l.Quantity
	}

	return 0
}
