package graph

import (
	"application/graph/model"
	"application/poker"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const (
	READER model.Role = "READER"
	WRITER model.Role = "WRITER"
)

type Resolver struct {
	Store poker.PlayerStore
	Role  model.Role
}

func Convert(player poker.Player) *model.Player {
	return &model.Player{
		Name: player.Name,
		Wins: player.Wins,
	}
}
