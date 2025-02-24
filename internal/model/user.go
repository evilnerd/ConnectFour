package model

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type User struct {
	Id    int64
	Name  string
	Email string
	Token string
}

func NewUser(name string, email string) User {
	return User{
		Name:  name,
		Email: email,
	}
}

func (u User) MakeToken() {
	u.Token = fmt.Sprintf("%x", md5.Sum([]byte(u.Name)))
}

func (u User) Is(other User) bool {
	return strings.EqualFold(u.Email, other.Email)
}

func (u User) Empty() bool {
	return u.Email == ""
}
func (u User) New() bool {
	return u.Id == 0
}
