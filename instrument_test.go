package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"gopkg.in/guregu/null.v4"
)

func TestInstruments(t *testing.T) {
	raw, err := os.ReadFile("instruments.json")
	if err != nil {
		t.Fatal(err)
	}

	var data struct {
		Table  string       `json:"table"`
		Action string       `json:"action"`
		Error  string       `json:"error"`
		Data   []Instrument `json:"data"`
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}

	it := NewInstrumentTable(data.Data)

	const btcPrice = 22755.17

	// Set btc price
	it.Update(Instrument{
		Symbol:    ".BXBT",
		MarkPrice: null.FloatFrom(btcPrice),
	})

	it.Update(Instrument{
		Symbol:    ".BXBTT30M",
		MarkPrice: null.FloatFrom(btcPrice),
	})

	it.Update(Instrument{
		Symbol:    ".BETHT",
		MarkPrice: null.FloatFrom(btcPrice),
	})

	// Test values
	table := []struct {
		Liq       RawLiquidation
		MarkPrice float64
		Quantity  float64
		Currency  string
		USD       float64
	}{
		// Linear contracts
		{
			Liq: RawLiquidation{
				LeavesQty: 1000000,
				Price:     22755.17,
				Symbol:    "XBTUSDTU23",
				Side:      "Buy",
			},
			Quantity: 1,
			Currency: "XBT",
			USD:      22758.5832755,
		},

		{
			Liq: RawLiquidation{
				LeavesQty: 29000,
				Price:     1636.46,
				Symbol:    "ETHUSDT",
				Side:      "Buy",
			},
			Quantity: 0.29,
			Currency: "ETH",
			USD:      474.6445860100001,
		},

		{
			Liq: RawLiquidation{
				LeavesQty: 50000,
				Price:     0.07221,
				Symbol:    "ETHH23",
				Side:      "Buy",
			},
			Quantity: 0.5,
			Currency: "ETH",
			USD:      821.5754128499999,
		},

		// Quanto contracts
		{
			Liq: RawLiquidation{
				LeavesQty: 2,
				Price:     0.086230,
				Symbol:    "DOGEUSD",
				Side:      "Buy",
			},
			Quantity: 2,
			Currency: "Cont",
			USD:      3.924,
		},

		{
			Liq: RawLiquidation{
				LeavesQty: 20,
				Price:     1592.55,
				Symbol:    "ETHUSD",
				Side:      "Buy",
			},
			Quantity: 20,
			Currency: "Cont",
			USD:      724.77,
		},

		// Inverse contracts
		{
			Liq: RawLiquidation{
				LeavesQty: 1000000,
				Price:     23245.5,
				Symbol:    "XBTUSD",
				Side:      "Buy",
			},
			Quantity: 1000000,
			Currency: "USD",
			USD:      1000000,
		},

		{
			Liq: RawLiquidation{
				LeavesQty: 1000000,
				Price:     23245.5,
				Symbol:    "XBTEUR",
				Side:      "Buy",
			},
			Quantity: 1000000,
			Currency: "EUR",
			USD:      1085600,
		},

		{
			Liq: RawLiquidation{
				LeavesQty: 7000,
				Price:     23649.48,
				Symbol:    "XBTH23",
				Side:      "Buy",
			},
			Quantity: 7000,
			Currency: "USD",
			USD:      7000,
		},
	}

	for _, v := range table {
		it.Update(Instrument{
			Symbol:    v.Liq.Symbol,
			MarkPrice: null.FloatFrom(v.Liq.Price),
		})

		l, err := it.Process(v.Liq)
		if err != nil {
			t.Fatal(err)
		}

		if int(l.TotalUSDValue)*100 != int(v.USD)*100 {
			t.Fatal("expected usd calculation", v.Liq, l.TotalUSDValue, v.USD)
		}
		if l.Currency != v.Currency {
			t.Fatal("expected currency", v.Liq, l.Currency, v.Currency)
		}
		if int(l.Quantity)*100 != int(v.Quantity)*100 {
			t.Fatal("expected quantity", v.Liq, l.Quantity, v.Quantity)
		}
		log.Println(l)
	}

	// {Price:1636.46, Quantity:0.29000000000000004, Currency:"ETH", TotalUSDValue:469.96240000000006, MinStep:0.01, MinTick:0.05}, Symbol:"ETHUSDT", Side:"Buy"}
	// {Price:0.07221, Quantity:0.5, Currency:"ETH", TotalUSDValue:810.28, MinStep:0.01, MinTick:1e-05}, Symbol:"ETHH23", Side:"Buy"}
	// {Price:0.08623, Quantity:2, Currency:"Cont", TotalUSDValue:3.9243566181999996, MinStep:1, MinTick:1e-05}, Symbol:"DOGEUSD", Side:"Buy"}
	// {Price:1592.55, Quantity:20, Currency:"Cont", TotalUSDValue:724.7749196699999, MinStep:1, MinTick:0.05}, Symbol:"ETHUSD", Side:"Buy"}
	// {Price:23245.5, Quantity:1e+06, Currency:"USD", TotalUSDValue:1e+06, MinStep:100, MinTick:0.5}, Symbol:"XBTUSD", Side:"Buy"}
	// {Price:23245.5, Quantity:1e+06, Currency:"EUR", TotalUSDValue:1.0856e+06, MinStep:100, MinTick:0.5}, Symbol:"XBTEUR", Side:"Buy"}
	// {Price:23649.48, Quantity:7000, Currency:"USD", TotalUSDValue:7000, MinStep:100, MinTick:0.5}, Symbol:"XBTH23", Side:"Buy"}

	cl := CombinedLiquidation{
		Symbol: "XBTUSDTU23",
		Side:   "Buy",
		Liquidations: []PriceQuantity{
			{Price: 22755.17, Quantity: 1, Currency: "XBT", TotalUSDValue: 22755.17, MinStep: 0.001, MinTick: 0.5},
			{Price: 22755.17, Quantity: 2, Currency: "XBT", TotalUSDValue: 22755.17, MinStep: 0.001, MinTick: 0.5},
		},
	}
	log.Println(cl)

	cl = CombinedLiquidation{
		Symbol: "XBTUSDTU23",
		Side:   "Buy",
		Liquidations: []PriceQuantity{
			{Price: 22755.0, Quantity: 1, Currency: "XBT", TotalUSDValue: 22755.17, MinStep: 0.001, MinTick: 0.5},
			{Price: 22755.1, Quantity: 2, Currency: "XBT", TotalUSDValue: 22755.18, MinStep: 0.001, MinTick: 0.5},
		},
	}
	log.Println(cl)
}

func TestInstruments2(t *testing.T) {
	raw, err := os.ReadFile("instruments2.json")
	if err != nil {
		t.Fatal(err)
	}

	var data struct {
		Table  string       `json:"table"`
		Action string       `json:"action"`
		Error  string       `json:"error"`
		Data   []Instrument `json:"data"`
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		t.Fatal(err)
	}

	it := NewInstrumentTable(data.Data)

	// Test values
	table := []struct {
		Liq       RawLiquidation
		MarkPrice float64
		Quantity  float64
		Currency  string
		USD       float64
	}{
		{
			Liq: RawLiquidation{
				LeavesQty: 4000,
				Price:     0.00635,
				Symbol:    "NEIROUSDT",
				Side:      "SELL",
			},
			Quantity: 4000,
			Currency: "NEIRO",
			USD:      25.409905999999996,
		},
	}

	for _, v := range table {
		it.Update(Instrument{
			Symbol:    v.Liq.Symbol,
			MarkPrice: null.FloatFrom(v.Liq.Price),
		})

		l, err := it.Process(v.Liq)
		if err != nil {
			t.Fatal(err)
		}

		if int(l.TotalUSDValue)*100 != int(v.USD)*100 {
			t.Fatal("expected usd calculation", v.Liq, l.TotalUSDValue, v.USD)
		}
		if l.Currency != v.Currency {
			t.Fatal("expected currency", v.Liq, l.Currency, v.Currency)
		}
		if int(l.Quantity)*100 != int(v.Quantity)*100 {
			t.Fatal("expected quantity", v.Liq, l.Quantity, v.Quantity)
		}
		log.Println(l)
	}
}
