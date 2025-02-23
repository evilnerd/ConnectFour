package models

import (
	"connectfour/internal/client/console/backend"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type AskNameModel struct {
	ConnectFourModel
	*State
	EmailText      textinput.Model
	PlayerPassword string
	PasswordText   textinput.Model
	ErrorMessage   string
	IsFailed       bool
}

func (s Step) String() string { return string(s) }

func NewAskNameModel(state *State) *AskNameModel {
	m := &AskNameModel{
		State:        state,
		EmailText:    textinput.New(),
		PasswordText: textinput.New(),
	}
	m.EmailText.Placeholder = "Your e-mail"
	m.EmailText.Focus()
	m.EmailText.CharLimit = 255
	m.EmailText.Width = 50
	m.PasswordText.Placeholder = "Your password"
	m.PasswordText.Focus()
	m.PasswordText.CharLimit = 100
	m.PasswordText.Width = 50
	return m
}

func (m AskNameModel) BreadCrumb() string {
	if m.PlayerName == "" {
		return "Name"
	}
	return fmt.Sprintf("Name (%s)", m.PlayerName)
}

func (m AskNameModel) Init() tea.Cmd {
	if m.wc.IsValid() && !m.wc.IsExpired() {
		m.PlayerName, m.PlayerEmail, _ = wc.Identify()
	}
	return nil
}

func (m AskNameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	// No need to ask for input if the WebClient already loaded
	// a valid JWT from disk.
	if m.wc.IsValid() && !m.wc.IsExpired() {
		m.PlayerName, m.PlayerEmail, _ = wc.Identify()
		return m.NextModel()
	}

	// If a form is active, make that one handle the key updates.
	switch msg := msg.(type) {

	case LoginMsg:
		if msg.isValid {
			m.IsFailed = false
			m.MustReauthenticate = false
			return m.NextModel()
		} else {
			m.IsFailed = true
			m.ErrorMessage = msg.errorMessage
		}
		break

	// Is it a key press?
	case tea.KeyMsg:

		switch msg.String() {
		case "esc", "ctrl+c":
			return m.PreviousModel()

		case "tab", "shift+tab":
			if !m.IsFailed {
				if m.EmailText.Focused() {
					m.EmailText.Blur()
					m.PasswordText.Focus()

				} else {
					m.PasswordText.Blur()
					m.EmailText.Focus()
				}
			}

		case "enter":
			if m.IsFailed {
				m.IsFailed = false
				m.ErrorMessage = ""
			} else if isValidEmail(m.EmailText.Value()) && isValidPassword(m.PasswordText.Value()) {
				m.PlayerEmail = m.EmailText.Value()
				m.PlayerPassword = m.PasswordText.Value()
				return m, m.Login()
			}
		}
	}

	var cmd tea.Cmd
	if m.EmailText.Focused() {
		m.EmailText, cmd = m.EmailText.Update(msg)
	} else {
		m.PasswordText, cmd = m.PasswordText.Update(msg)
	}
	return m, cmd
}

func isValidEmail(val string) bool {
	l := len(strings.TrimSpace(val))
	return l > 2 && l <= 255
}

func isValidPassword(val string) bool {
	l := len(strings.TrimSpace(val))
	return l > 2 && l <= 100
}

func (m AskNameModel) View() string {
	var view string

	if m.wc.IsValid() {
		view = lipgloss.JoinVertical(lipgloss.Left,
			styles.Label.Render("Welcome back "),
			styles.Value.Render(m.PlayerName),
		)
	} else if m.IsFailed {
		view = lipgloss.JoinVertical(lipgloss.Left,
			styles.Label.Render("There was a problem"),
			styles.Value.Render(m.ErrorMessage),
		)
	} else {
		view = lipgloss.JoinVertical(lipgloss.Left,
			styles.Description.Render("Enter your e-mail to uniquely identify you"),
			m.EmailText.View(),
			styles.Description.Render("Enter the password you set"),
			m.PasswordText.View(),
		)
	}

	return m.CommonView(view)
}

func (m AskNameModel) Login() tea.Cmd {
	return func() tea.Msg {
		err := backend.Login(m.wc, m.PlayerEmail, m.PlayerPassword)
		out := LoginMsg{}
		if err != nil {
			out.errorMessage = err.Error()
			out.isValid = false
		} else {
			out.isValid = true
		}
		return out
	}
}
