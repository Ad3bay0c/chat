package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)
type Server struct {
	Groups map[string]*Group
	Users map[net.Addr]*User
	Commands chan *Command
}

var s = &Server{
	Groups: make(map[string]*Group),
	Users: make(map[net.Addr]*User),
	Commands: make(chan *Command),
}
//
//var s = &Server{
//	Chats: make(map[string]*Chat),
//	Users: make(map[net.Addr]*User),
//	Commands: make(chan string),
//}
//func StartServer() {
//	lis, err := net.Listen("tcp", ":8080")
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Listening on %s\n", lis.Addr())
//	for {
//		conn, err := lis.Accept()
//		if err != nil {
//			panic(err)
//		}
//		fmt.Printf("A new conection: %s\n", conn.RemoteAddr().String())
//
//		go s.handleInput(conn)
//	}
//}
//
//func (s *Server) handleInput(conn net.Conn) {
//	user := &User{
//		conn: conn,
//		userName: "unknown",
//	}
//	s.Users[conn.RemoteAddr()] = user
//
//	for {
//		message, err := bufio.NewReader(conn).ReadString('\n')
//		if err != nil {
//			conn.Write([]byte("Error reading from connection"))
//			return
//		}
//		message = strings.Trim(message, "\n") + " : from server"
//		conn.Write([]byte(message))
//	}
//}

func StartServer() {
	go s.readCommands()
	// socket => server and client
	// server => it accepts request (read) from the client and sends (write) back response
	// client => it sends (write) request to the server and receives (read) response
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	log.Println("Listening on", lis.Addr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("A new conection: %s\n", conn.RemoteAddr().String())
		user := &User{
			conn: conn,
			userName: "unknown",
		}
		s.Users[conn.RemoteAddr()] = user
		go s.readInput(user)
		//log.Println(message)
		//conn.Write([]byte(fmt.Sprintf("Server: %s\n", message)))
	}

}

// main thread
// readInput thread
func (s *Server) readInput(user *User) {
	for {
		message, err := bufio.NewReader(user.conn).ReadString('\n')
		if err != nil && err == io.EOF {
			log.Println(err)
			log.Println("Client disconnected: ", user.conn.RemoteAddr())
			user.conn.Close()
			return
		}
		message = strings.Trim(message, "\n")
		commands := strings.Split(message, " ")
		command := commands[0]
		switch command {
		case "join":
			s.Commands <- &Command{
				Name: COMMAND_JOIN,
				User: user,
				Args: commands[1:],
			}
		case "send":
			s.Commands <- &Command{
				Name: COMMAND_SEND,
				User: user,
				Args: commands[1:],
			}
		}
		//user.group.broadcast(user.conn, fmt.Sprintf("Server: %s\n", message))
	}
}

func (s *Server) readCommands() {
	for command := range s.Commands {
		switch command.Name {
		case COMMAND_JOIN:
			s.JoinGroup(command.User, command.Args[0])
		case COMMAND_SEND:
			s.SendMessage(command.User, command.Args)
		}
	}
}
func (s *Server) JoinGroup(user *User, groupName string) {
	grp, ok := s.Groups[groupName]
	if !ok {
		grp = &Group{
			name: groupName,
			members: make(map[net.Addr]*User),
		}
		s.Groups[groupName] = grp
	}
	user.group = grp
	grp.members[user.conn.RemoteAddr()] = user
}

func (s *Server) SendMessage(user *User, message []string) {
	newMessage := strings.Join(message[:], " ")
	user.group.broadcast(user.conn, fmt.Sprintf("%s: %s\n", user.userName, newMessage))
}