package console

type Option struct {
	key         string
	title       string
	description string
}

func NewOption(key string, title string, description string) Option {
	return Option{
		key,
		title,
		description,
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

func (i Option) Key() string { return i.key }
