package main

type User struct {
	Email          string // also serves as primary key?
	Name           string
	HashedPassword string
	SessionToken string
	CSRFToken string
}
