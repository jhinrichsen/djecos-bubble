package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestNewDeck(t *testing.T) {
	d := NewDeck()
	want := 32
	got := len(d)
	if want != got {
		t.Fatalf("want %d but got %d\n", want, got)
	}
}

func TestShuffle(t *testing.T) {
	d := NewDeck()
	Shuffle(&d)
	// How do you test randomness?
	/*
		for _, card := range d {
			log.Printf("%+v \n", card)
		}
	*/
}

func TestGameOver(t *testing.T) {
	var o Octopus
	for i := range o {
		o[i] = true
	}
	want := true
	got := GameOver(o)
	if want != got {
		t.Fatalf("want %t but got %t\n", want, got)
	}
}

func TestTwoPlayerGame(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	d := NewDeck()
	Shuffle(&d)
	var strategies = []DrawAgain{
		AlwaysDraw,
		NeverDraw,
	}

	NewGame(d, strategies)
}

func TestMultiTwoPlayerGames(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var wins [2]int
	totalTurns := 0
	for i := 0; i < 1000; i++ {
		d := NewDeck()
		Shuffle(&d)
		var strategies = []DrawAgain{
			AlwaysDraw,
			NeverDraw,
		}
		winner, turns := NewGame(d, strategies)
		wins[winner]++
		totalTurns += turns
	}
	fmt.Printf("player 0: %d wins, player 1: %d wins, %d turns avg\n",
		wins[0], wins[1], totalTurns/1000)
}

func TestMultiFourPlayerGames(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	var wins [4]int
	totalTurns := 0
	for i := 0; i < 1000; i++ {
		d := NewDeck()
		Shuffle(&d)
		var strategies = []DrawAgain{
			AlwaysDraw,
			NeverDraw,
			AlwaysDraw,
			NeverDraw,
		}
		winner, turns := NewGame(d, strategies)
		wins[winner]++
		totalTurns += turns
	}
	fmt.Printf("player 0: %d wins, player 1: %d wins\n", wins[0], wins[1])
	fmt.Printf("player 2: %d wins, player 3: %d wins\n", wins[2], wins[3])
	fmt.Printf("%d turns avg\n", totalTurns/1000)
}
