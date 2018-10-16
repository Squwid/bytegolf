package main

import "testing"

func TestNewGame(t *testing.T) {
	game, err := NewGame("Bens Game", "", "medium", 9, 5)
	if err != nil {
		t.Logf("error creating game: %v\n", err)
		t.Fail()
	}
	if game.Name != "Bens Game" {
		t.Fail()
	}
}

func TestAddPlayers(t *testing.T) {
	game, err := NewGame("Bens Game", "", "medium", 9, 5)
	if err != nil {
		t.Logf("error creating game: %v\n", err)
		t.Fail()
	}
}
