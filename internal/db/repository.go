package db

import (
	"connectfour/internal/model"
)

type UserRepository interface {
	Create(u model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

type GameRepository interface {
	Save(game model.Game) bool
	Fetch(key string) (model.Game, error)
	List(userId int64, status string) ([]model.Game, error)
}
