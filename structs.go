package main

import "time"

type User struct {
	Email          string // also serves as primary key?
	Name           string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

type BusinessApplication struct {
	ApplicationID   string
	UserEmail       string
	ApplicationType string
	BusinessName    string
	BusinessType    string
	Description     string
	Address         string
	DateCreated     time.Time
}
