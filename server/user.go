package server

import (
	"crypto/md5"
	"fmt"
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
