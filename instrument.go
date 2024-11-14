package main

import (
	"errors"
	"math"
	"sort"
	"strings"

	"gopkg.in/guregu/null.v4"
)

// InstrumentType as defined by Bitmex.
// Ref: https://www.bitmex.com/api/explorer/#!/Instrument/Instrument_get
type InstrumentType string

// InstrumentType definitions
const (
	ITPerpetualContracts             InstrumentType = "FFWCSX" // Perpetual Contracts
	ITPerpetualContractsFXUnderliers InstrumentType = "FFWCSF" // Perpetual Contracts (FX underliers)
	ITSpot                           InstrumentType = "IFXXXP" // Spot
	ITFutures                        InstrumentType = "FFCCSX" // Futures

	ITBMBasketIndex     InstrumentType = "MRBXXX" // BitMEX Basket Index
	ITBMCryptoIndex     InstrumentType = "MRCXXX" // BitMEX Crypto Index
	ITBMFXIndex         InstrumentType = "MRFXXX" // Bitmex FX Index
	ITBMPremiumIndex    InstrumentType = "MRRXXX" // Bitmex Premium Index
	ITBMVolatilityIndex InstrumentType = "MRIXXX" // Bitmex Volatility Index
)

// Instrument an instrument on Bitmex.
type Instrument struct {
	Symbol Symbol         `json:"symbol"`
	Type   InstrumentType `json:"typ"`

	AskPrice   null.Float `json:"askPrice"`
	BidPrice   null.Float `json:"bidPrice"`
	LastPrice  null.Float `json:"lastPrice"`
	MarkPrice  null.Float `json:"markPrice"`
	TickSize   null.Float `json:"tickSize"`
	LotSize    null.Float `json:"lotSize"`
	Multiplier null.Float `json:"multiplier"`

	IsInverse bool `json:"isInverse"`
	IsQuanto  bool `json:"isQuanto"`

	ReferenceSymbol                Symbol     `json:"referenceSymbol"`
	PositionCurrency               string     `json:"positionCurrency"`
	QuoteCurrency                  string     `json:"quoteCurrency"`
	SettleCurrency                 string     `json:"settlCurrency"`
	QuoteToSettleMultiplier        null.Float `json:"quoteToSettleMultiplier"`
	Underlying                     string     `json:"underlying"`
	UnderlyingSymbol               string     `json:"underlyingSymbol"`
	UnderlyingToPositionMultiplier null.Float `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   null.Float `json:"underlyingToSettleMultiplier"`
}

// InstrumentTable of all instruments, used for queries and contract value calculations.
type InstrumentTable struct {
	insts map[Symbol]Instrument
}

// NewInstrumentTable creates a new instrument table.
func NewInstrumentTable(insts []Instrument) *InstrumentTable {
	table := InstrumentTable{
		insts: make(map[Symbol]Instrument),
	}

	for _, v := range insts {
		table.insts[v.Symbol] = v
	}

	return &table
}

// Update the value of an instrument.
func (it *InstrumentTable) Update(update Instrument) {
	inst, ok := it.insts[update.Symbol]
	if !ok {
		it.insts[update.Symbol] = update
		return
	}

	if update.AskPrice.Valid {
		inst.AskPrice = update.AskPrice
	}
	if update.BidPrice.Valid {
		inst.BidPrice = update.BidPrice
	}
	if update.LastPrice.Valid {
		inst.LastPrice = update.LastPrice
	}
	if update.MarkPrice.Valid {
		inst.MarkPrice = update.MarkPrice
	}
	if update.TickSize.Valid {
		inst.TickSize = update.TickSize
	}
	if update.Multiplier.Valid {
		inst.Multiplier = update.Multiplier
	}
	if update.QuoteToSettleMultiplier.Valid {
		inst.QuoteToSettleMultiplier = update.QuoteToSettleMultiplier
	}
	if update.UnderlyingToPositionMultiplier.Valid {
		inst.UnderlyingToPositionMultiplier = update.UnderlyingToPositionMultiplier
	}
	if update.UnderlyingToSettleMultiplier.Valid {
		inst.UnderlyingToSettleMultiplier = update.UnderlyingToSettleMultiplier
	}

	it.insts[update.Symbol] = inst
}

// PriceUSD returns price of 1 unit of currency in USD.
// Returns 0 if not found.
func (it *InstrumentTable) PriceUSD(currency string) float64 {
	if currency == "USD" {
		return 1
	}

	if currency == "XBt" {
		return it.insts[".BXBT"].MarkPrice.Float64 / 100000000
	} else if currency == "XBT" {
		return it.insts[".BXBT"].MarkPrice.Float64
	}

	if currency == "USDt" {
		return it.insts[".BUSDT"].MarkPrice.Float64 / 100000000
	} else if currency == "USDT" {
		return it.insts[".BUSDT"].MarkPrice.Float64
	}

	var found []Symbol

	// Use the instruments to look up price of the position currency to USD
	for k, v := range it.insts {
		switch v.Type {
		case ITBMBasketIndex, ITBMCryptoIndex, ITBMFXIndex:
			if strings.ToUpper(v.QuoteCurrency) == "USD" && strings.ToUpper(v.Underlying) == currency {
				if v.MarkPrice.Valid {
					found = append(found, k)
				}
			}
		}
	}

	// Use the symbol with the shortest name
	sort.Slice(found, func(i, j int) bool {
		return len(found[i]) < len(found[j])
	})

	if len(found) != 0 {
		return it.insts[found[0]].MarkPrice.Float64
	}

	return 0
}

func (it *InstrumentTable) USDLookup(symbol Symbol) (res float64) {
	inst, ok := it.insts[symbol]
	if !ok {
		return 0
	}

	if inst.QuoteCurrency == "USD" {
		return 1
	} else if inst.QuoteCurrency == "XBt" {
		return it.insts[".BXBT"].MarkPrice.Float64 / 100000000
	} else if inst.QuoteCurrency == "XBT" {
		return it.insts[".BXBT"].MarkPrice.Float64
	} else if inst.QuoteCurrency == "USDt" {
		return it.insts[".BUSDT"].MarkPrice.Float64 / 100000000
	} else if inst.QuoteCurrency == "USDT" {
		return it.insts[".BUSDT"].MarkPrice.Float64
	}

	if inst.ReferenceSymbol == symbol {
		return 0
	}

	return it.USDLookup(inst.ReferenceSymbol) * inst.MarkPrice.Float64
}

// Process calculates the liquidated position's display quantity, units and USD value.
func (it *InstrumentTable) Process(rl RawLiquidation) (Liquidation, error) {
	inst, ok := it.insts[rl.Symbol]
	if !ok {
		return Liquidation{}, errors.New("instrument not found")
	}

	currency := inst.PositionCurrency
	if inst.PositionCurrency == "" {
		currency = "Cont"
	}

	pq := PriceQuantity{
		Price:    rl.Price,
		Currency: currency,
		MinStep:  inst.MinStep(),
		MinTick:  inst.TickSize.Float64,
	}

	if inst.IsInverse {
		pq.Quantity = rl.LeavesQty
		pq.TotalUSDValue = rl.LeavesQty * it.PriceUSD(inst.PositionCurrency)
	} else if inst.IsQuanto {
		pq.Quantity = rl.LeavesQty
		pq.TotalUSDValue = inst.NotionalValue() * rl.LeavesQty * it.PriceUSD(inst.SettleCurrency)
	} else {
		pq.Quantity = rl.LeavesQty * inst.NotionalValue()
		pq.TotalUSDValue = inst.NotionalValue() * rl.LeavesQty * it.USDLookup(inst.ReferenceSymbol) * inst.MarkPrice.Float64
	}

	// log.Println(inst.Symbol)
	// log.Println("Calculated position:", pq.Quantity, pq.Currency)
	// log.Println("Calculated value: $", pq.TotalUSDValue)
	// log.Println("Calculated size: $", pq.Quantity)
	// log.Println()

	return Liquidation{
		PriceQuantity: pq,
		Symbol:        rl.Symbol,
		Side:          rl.Side,
	}, nil
}

// ContractValueXBT ...
func (in Instrument) ContractValueXBT() float64 {
	var c float64
	if int64(in.Multiplier.Float64) == 1 {
		c = in.LotSize.Float64 * in.MarkPrice.Float64
	} else {
		if in.Multiplier.Float64 > 0 {
			c = in.LotSize.Float64 * in.Multiplier.Float64 * in.MarkPrice.Float64
		} else {
			c = in.LotSize.Float64 * in.Multiplier.Float64 / in.MarkPrice.Float64
		}
	}

	return math.Abs(c / in.LotSize.Float64)
}

// ContractValue ...
func (in Instrument) ContractValue() float64 {
	if in.IsQuanto {
		return in.ContractValueXBT()
	}

	if in.IsInverse {
		return math.Abs(in.Multiplier.Float64 / in.UnderlyingToSettleMultiplier.Float64)
	}

	return 1
}

// NotionalValue denominated in the settlCurrency.
func (in Instrument) NotionalValue() float64 {
	if in.IsInverse {
		return in.ContractValue() * in.Multiplier.Float64 / in.UnderlyingToSettleMultiplier.Float64
	}

	if in.IsQuanto {
		return in.ContractValue()
	}

	if in.UnderlyingToPositionMultiplier.Valid {
		return in.ContractValue() / in.UnderlyingToPositionMultiplier.Float64
	}

	return in.ContractValue()
}

// MinStep returns the minimum amount a position can be entered.
func (in Instrument) MinStep() float64 {
	if in.IsInverse {
		return in.LotSize.Float64 * in.Multiplier.Float64 / in.UnderlyingToSettleMultiplier.Float64
	}

	if in.IsQuanto {
		return in.LotSize.Float64
	}

	if in.UnderlyingToPositionMultiplier.Valid {
		return in.LotSize.Float64 / in.UnderlyingToPositionMultiplier.Float64
	}

	return in.LotSize.Float64
}
