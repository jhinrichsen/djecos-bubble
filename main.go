package main

import (
	"log"
	"math/rand"
)

// Style represents the visual appearance of a card
type Style int

const (
	balloons Style = iota
	fireworks
	stars
	relaxed
	favourite
)

// Action determines what to do
type Action int

const (
	pattern Action = iota
	staringEyes
	tentacle
	staringOctopus
)

// Card has a style and an action
type Card struct {
	Style  Style
	Action Action
}

// Deck is the complete set of all cards
type Deck [32]Card

// NewDeck returns an unshuffled set of all cards
func NewDeck() Deck {
	return Deck{

		// 4x5 cards (four for each Style)
		{balloons, pattern}, {balloons, pattern},
		{balloons, pattern}, {balloons, pattern},
		{fireworks, pattern}, {fireworks, pattern},
		{fireworks, pattern}, {fireworks, pattern},
		{stars, pattern}, {stars, pattern},
		{stars, pattern}, {stars, pattern},
		{relaxed, pattern}, {relaxed, pattern},
		{relaxed, pattern}, {relaxed, pattern},
		{favourite, pattern}, {favourite, pattern},
		{favourite, pattern}, {favourite, pattern},

		// 5 staringEyes, one for each Style
		{balloons, staringEyes},
		{fireworks, staringEyes},
		{stars, staringEyes},
		{relaxed, staringEyes},
		{favourite, staringEyes},

		// 5 tentacles, one for each Style
		{balloons, tentacle},
		{fireworks, tentacle},
		{stars, tentacle},
		{relaxed, tentacle},
		{favourite, tentacle},

		// 2 staring octopus, independent of associated style
		{balloons /* irrelevant */, staringOctopus},
		{balloons /* irrelevant */, staringOctopus},
	}
}

// Shuffle in place
// TODO how much faster is in place vs. copy?
func Shuffle(d *Deck) {
	// As of Go 1.10
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

// Octopus has 5 arms
type Octopus [5]bool

// DrawAgain either draws (true) or passes (false)
type DrawAgain func() bool

// NewGame returns number of turns (which is not number of draws)
// each player is identified by a strategy
func NewGame(deck Deck, strategies []DrawAgain) (player, turns int) {
	// Each player has one octopus
	octopus := make([]Octopus, len(strategies))
	drawn := 0
	var canDraw bool
	for {
		turns++

	draw:
		card := deck[drawn]
		// Deck exhausted?
		drawn++
		if drawn == len(deck) {
			drawn = 0
			log.Printf("empty deck, re-shuffling\n")
			Shuffle(&deck)
		}

		if card.Action == staringOctopus {
			// Player loses all fishes
			log.Printf("Player %d loses all fishes\n", player)
			for i := range octopus[player] {
				octopus[player][i] = false
			}
			canDraw = false
		} else if card.Action == staringEyes {
			hasFish := &octopus[player][card.Style]
			if *hasFish {
				// Player loses fish on given style
				*hasFish = false
				canDraw = false
			} else {
				canDraw = true
			}
		} else if card.Action == tentacle {
			hasFish := &octopus[player][card.Style]
			if *hasFish {
				canDraw = true
			} else {
				// Steal fish from another player
				Steal(player, card.Style, octopus)
				canDraw = false
			}
		} else {
			log.Printf("player %d picks a %v\n", player, card.Style)
			// Plain pattern
			octopus[player][card.Style] = true
			canDraw = true
		}
		log.Printf("player %d has %d fishes\n", player, Fishes(octopus[player]))
		if GameOver(octopus[player]) {
			log.Printf("player %d wins in %d turns\n", player, turns)
			return
		}

		// Draw or pass?
		if canDraw {
			log.Printf("player %d can draw", player)
			if strategies[player]() {
				log.Printf("player %d chooses to draw", player)
				// Look mom - not even structured programming in 2018
				goto draw
			} else {
				log.Println("player chooses to pass")
			}
		}

		// Next player
		player = (player + 1) % len(strategies)
	}
}

// GameOver returns true if all arms are filled by a fish
func GameOver(o Octopus) bool {
	return Fishes(o) == len(o)
}

// Fishes returns number of arms that have a fish
func Fishes(o Octopus) int {
	n := 0
	for _, arm := range o {
		if arm {
			n++
		}
	}
	return n
}

// AlwaysDraw keeps drawing new cards as long as possible
func AlwaysDraw() bool {
	return true
}

// NeverDraw always passes
func NeverDraw() bool {
	return false
}

// Steal a fish from the player with the most fishes
func Steal(player int, style Style, octopus []Octopus) {
	stealFromPlayer := -1
	maxFishes := 0
	log.Printf("player %d wants to steal a fish style %d\n", player, style)
	for i := range octopus {
		// Don't steal from ourself
		if i == player {
			continue
		}

		// Player has a fish for given style?
		if octopus[i][style] {
			log.Printf("player %d has a matching fish\n", i)
			n := Fishes(octopus[i])
			if n > maxFishes {
				maxFishes = n
				stealFromPlayer = i
			}
		}
	}
	if stealFromPlayer == -1 {
		log.Printf("no player has a matching fish, bad luck for player %d\n", player)
	} else {
		log.Printf("player %d steals from player %d\n", player, stealFromPlayer)
		octopus[stealFromPlayer][style] = false
		octopus[player][style] = true
	}
}
