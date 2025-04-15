package models

type User struct {
	Id       int
	Email    string
	PassHash []byte
}