package server

type Command struct {
	Name int
	Args []string
	User *User
}

const (
	COMMAND_JOIN = iota + 1
	COMMAND_SEND
)
