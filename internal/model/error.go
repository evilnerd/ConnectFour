package model

import "fmt"

type UnknownGameError struct {
	key string
}

func NewUnknownGameError(key string) UnknownGameError {
	return UnknownGameError{
		key,
	}
}

func (e UnknownGameError) Error() string {
	return fmt.Sprintf("the requested game key '%s' was not found", e.key)
}
