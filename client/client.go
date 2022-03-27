package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags| log.Lshortfile | log.Llongfile)
	//
	//conn, err := net.Dial("tcp", "localhost:8080")
	//if err != nil {
	//	panic(err)
	//}
	//for {
	//	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	_, err = conn.Write([]byte(msg))
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	go readFromServer(conn)
	//}
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	log.Println(conn.LocalAddr())
	go readFromServer(conn)
	for {
		msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	}


}

func readFromServer(conn net.Conn) {
	for {
		buf := make([]byte, 1024*1024)
		n, err := conn.Read(buf)
		if err != nil && err == io.EOF {
			return
		}
		message := strings.Trim(string(buf[:n]), "\n")
		log.Println(string(message))
	}
}

