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
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(500000)),
			Symbol:   symbols[i%len(symbols)],
			Side:     "Buy",
		}).String()

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
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(500000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		if err := s.Save(); err != nil {
			t.Fatal(err)
		}

		log.Println(result)
		verify(result, t)
	}

	time.Sleep(22 * time.Second)

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(500000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		log.Println(result)
		verify(result, t)
	}

	time.Sleep(22 * time.Second)

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(500000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		time.Sleep(3 * time.Second)

		log.Println(result)
		verify(result, t)
	}

	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
}

func Test10m(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	result := s.Decorate(Liquidation{
		Price:    10000,
		Quantity: 10000000,
		Symbol:   "XBTUSD",
		Side:     "Buy",
	}).String()

	log.Println(result)
	verify(result, t)

	result = s.Decorate(Liquidation{
		Price:    10000,
		Quantity: 100000000,
		Symbol:   "XBTUSD",
		Side:     "Buy",
	}).String()

	log.Println(result)
	verify(result, t)

	result = s.Decorate(Liquidation{
		Price:    10000,
		Quantity: 1000000000,
		Symbol:   "XBTUSD",
		Side:     "Buy",
	}).String()

	log.Println(result)
	verify(result, t)
}

func TestSaveFile(t *testing.T) {
	s, err := NewState()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		result := s.Decorate(Liquidation{
			Price:    float64(rand.Intn(i + 1)),
			Quantity: int64(rand.Intn(500000)),
			Symbol:   "BTCUSD",
			Side:     "Buy",
		}).String()

		if err := s.Save(); err != nil {
			t.Fatal(err)
		}

		log.Println(result)
		verify(result, t)
	}
}
