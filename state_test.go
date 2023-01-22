package main

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

// TODO: fix these tests so they work without inspection

func verify(result string, t *testing.T) {
	if len([]rune(result)) > twitterLengthLimit {
		t.Fatal("longer than the length limit")
	}
}

func TestSymbolLiquidator(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	liqChan := make(chan Liquidation)
	tweetChan := make(chan preparedTweet)
	go symbolLiquidator(s, liqChan, tweetChan)

	go func() {
		for result := range tweetChan {
			// It is a lot easier to test by inspection
			log.Println(result)
			verify(result.status, t)
		}
	}()

	// Generate liquidations, expect no panics or errors
	for i := 0; i < 10; i++ {
		l := Liquidation{
			PriceQuantity: PriceQuantity{
				Price:    float64(5000 + rand.Intn(i+1)),
				Quantity: float64(rand.Intn(10)) + 1,
				Currency: "USD",
			},
			Symbol: "XBTUSD",
			Side:   "Buy",
		}

		liqChan <- l
	}
	close(liqChan)
}

func TestStateSimple(t *testing.T) {
	symbols := map[int]Symbol{
		0: "XBTUSD",
		1: "XBTZ16",
		2: "XBJ24H",
	}

	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	// Generate 100k liquidations, expect no panics or errors
	for i := 0; i < 100000; i++ {
		l := Liquidation{
			PriceQuantity: PriceQuantity{
				Price:    float64(rand.Intn(i + 1)),
				Quantity: float64(rand.Intn(500000)),
				Currency: "USD",
			},
			Symbol: symbols[i%len(symbols)],
			Side:   "Buy",
		}
		result := s.Decorate(l.ToCombined()).Apply(l.String())

		// It is a lot easier to test by inspection
		log.Println(result)
		verify(result, t)
	}
}

func TestStreaks(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		l := Liquidation{
			PriceQuantity: PriceQuantity{
				Price:    float64(rand.Intn(i + 1)),
				Quantity: float64(rand.Intn(500000)),
				Currency: "USD",
			},
			Symbol: "BTCUSD",
			Side:   "Buy",
		}

		result := s.Decorate(l.ToCombined()).Apply(l.String())
		log.Println(result)
		verify(result, t)
	}

	time.Sleep(22 * time.Second)

	for i := 0; i < 10; i++ {
		l := Liquidation{
			PriceQuantity: PriceQuantity{
				Price:    float64(rand.Intn(i + 1)),
				Quantity: float64(rand.Intn(500000)),
				Currency: "USD",
			},
			Symbol: "BTCUSD",
			Side:   "Buy",
		}
		result := s.Decorate(l.ToCombined()).Apply(l.String())

		log.Println(result)
		verify(result, t)
	}

	time.Sleep(22 * time.Second)

	for i := 0; i < 10; i++ {
		l := Liquidation{
			PriceQuantity: PriceQuantity{
				Price:    float64(rand.Intn(i + 1)),
				Quantity: float64(rand.Intn(500000)),
				Currency: "USD",
			},
			Symbol: "BTCUSD",
			Side:   "Buy",
		}
		result := s.Decorate(l.ToCombined()).Apply(l.String())

		time.Sleep(3 * time.Second)

		log.Println(result)
		verify(result, t)
	}
}

func Test10m(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	l := Liquidation{
		PriceQuantity: PriceQuantity{
			Price:    10000,
			Quantity: 10000000,
			Currency: "USD",
		},
		Symbol: "XBTUSD",
		Side:   "Buy",
	}

	result := s.Decorate(l.ToCombined()).Apply(l.String())

	log.Println(result)
	verify(result, t)

	l = Liquidation{
		PriceQuantity: PriceQuantity{
			Price:    10000,
			Quantity: 100000000,
			Currency: "USD",
		},
		Symbol: "XBTUSD",
		Side:   "Buy",
	}

	result = s.Decorate(l.ToCombined()).Apply(l.String())

	log.Println(result)
	verify(result, t)

	l = Liquidation{
		PriceQuantity: PriceQuantity{
			Price:    10000,
			Quantity: 1000000000,
			Currency: "USD",
		},
		Symbol: "XBTUSD",
		Side:   "Buy",
	}

	result = s.Decorate(l.ToCombined()).Apply(l.String())

	log.Println(result)
	verify(result, t)
}
