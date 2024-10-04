package console

type Option struct {
	title       string
	description string
}

func NewOption(title string, description string) Option {
	return Option{
		title, description,
	}
}

func (i Option) Title() string {
	return i.title
}

func (i Option) Description() string {
	return i.description
}

func (i Option) FilterValue() string {
	return i.title
}
