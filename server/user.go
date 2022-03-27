package server

import (
	"fmt"
	"net"
)

type User struct {
	userName string
	conn net.Conn
	group *Group
}

func (u *User) WriteMessage(message string) {
	u.conn.Write([]byte(fmt.Sprintf("%s-> %s\n", u.userName, message)))
}