package console

import "github.com/charmbracelet/lipgloss"

type AppStyles struct {
	// The Normal state.
	AppTitle            lipgloss.Style
	Label               lipgloss.Style
	Value               lipgloss.Style
	Header              lipgloss.Style
	Subdued             lipgloss.Style
	BreadCrumb          lipgloss.Style
	BreadCrumbSeparator lipgloss.Style
	Page                lipgloss.Style
	Description         lipgloss.Style
	Error               lipgloss.Style
}

// NewAppStyles returns style definitions for the app. So that we can refer to styles by name and manage them
// in a single place.
func NewAppStyles() (s AppStyles) {
	s.AppTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(1, 0, 2)

	s.Label = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 1, 0, 0)

	s.Value = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	s.Header = s.Value.Copy().
		Bold(true)

	s.Subdued = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.BreadCrumb = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#4D4D4D", Dark: "#9D9DAD"}).
		Align(lipgloss.Left)

	s.BreadCrumbSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAA")).
		Bold(true).
		Align(lipgloss.Left)

	s.Page = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).Padding(2)

	s.Description = s.Label.Copy().
		Padding(0, 0, 1, 0).
		Background(lipgloss.AdaptiveColor{Light: "#222", Dark: "DDD"})

	s.Error = s.Description

	return s
}
