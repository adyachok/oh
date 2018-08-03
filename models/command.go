package models

import (
	"time"
)

type Command struct {
	Created 		time.Time
	CommandType 	string  // "upload", "sync"
	State 			string	// "initiated", "in progress", "finised", "error"
	Filename 		string
}

func NewCommand(ct, filename string) *Command {
	return &Command {
		Created: time.Now(),
		CommandType:  ct,
		State: "initiated",
		Filename: filename,
	}
}