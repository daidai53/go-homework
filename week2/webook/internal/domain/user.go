// Copyright@daidai53 2023
package domain

type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	Phone    string
	Birthday string
	AboutMe  string
}
