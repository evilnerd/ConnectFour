package client

import (
	"errors"
	"github.com/charmbracelet/huh"
	"strings"
)

const (
	NewGameKey    = "new_game"
	PlayernameKey = "player_name"
)

func NewGetTheNameForm(model StartModel) *huh.Form {

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Play ConnectFour").
				Description("Let's get some information to either start or join a ConnectFour game").
				Next(true),
			huh.NewInput().
				Key(PlayernameKey).
				Value(&model.LocalPlayerName).
				Title("What's your name?").
				Placeholder("your name").
				Validate(
					func(s string) error {
						if strings.TrimSpace(s) == "" {
							return errors.New("you really should enter a name")
						}
						return nil
					}).
				Description("This is the name the other player(s) will see.").
				CharLimit(50),
			huh.NewConfirm().
				Key(NewGameKey).
				Value(&model.IsNewGame).
				Title("Start a new game?").
				Description("Do you want to start a new game? If not, you can join an existing game.").
				Affirmative("Yes, start a new game").
				Negative("No, look for an existing game"),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Join a private game?").
				Description("Do you want to enter the key of a private game you were invited to, or look for a public game to join?").
				Affirmative("Yes, join a private game").
				Negative("No, look for a public game to join").
				Value(&model.IsPrivateGame),
		).WithHideFunc(func() bool { return model.IsNewGame }),
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the game key").
				Description("Enter the key your opponent gave you to join their game.")),
	).WithAccessible(true)

	return form

}
