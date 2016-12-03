package main

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

// TODO: fix these tests so they work without inspection

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
	for i := 0; i < 100; i++ {
		_ = s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(100000)),
			Symbol:   symbols[i%len(symbols)],
			Side:     "Buy",
		}).String()

		// It is a lot easier to test by inspection
		// log.Println(result)
	}
}

func TestStreaks(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(100000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		if err := s.Save(); err != nil {
			t.Fatal(err)
		}

		log.Println(result)
	}

	time.Sleep(12 * time.Second)

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(100000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		log.Println(result)
	}

	time.Sleep(12 * time.Second)

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(100000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		time.Sleep(3 * time.Second)

		log.Println(result)
	}

	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
}

func TestSaveFile(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(100000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		if err := s.Save(); err != nil {
			t.Fatal(err)
		}

		log.Println(result)
	}
}
