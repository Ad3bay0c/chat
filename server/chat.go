package server

import (
	"net"
)

type Group struct {
	name string
	members map[net.Addr]*User
}

func (g *Group) broadcast(conn net.Conn, message string) {
	for userAddr, user := range g.members {
		if userAddr != conn.RemoteAddr() {
			user.WriteMessage(message)
		}
	}
}
