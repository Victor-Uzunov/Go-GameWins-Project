package graph

import (
	"application/graph/model"
	"application/poker"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Store poker.PlayerStore
}

func Convert(player poker.Player) *model.Player {
	return &model.Player{
		Name: player.Name,
		Wins: player.Wins,
	}
}
