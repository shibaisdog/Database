package main

import (
	"cloud-logic/database/memory"
	"cloud-logic/server/router"
	"log"
	"net"
	"os"
)

func Handle(db *memory.Storage, conn net.Conn) {
	defer conn.Close()
	client := router.Client{
		Login: false,
	}
	router.LoggerConnet(conn)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				router.LoggerDisconnetWithMessage(conn, "Timeout")
				break
			} else if err.Error() == "EOF" {
				router.LoggerDisconnet(conn)
				break
			}
			router.LoggerDisconnetWithMessage(conn, err.Error())
			break
		}
		router.Router(db, &client, conn, buffer[:n])
	}
}

func main() {
	DB := memory.NewStorage(os.Args[1])
	listener, err := net.Listen("tcp", os.Args[2])
	if err != nil {
		log.Fatalln("Error starting TCP server:", err)
	}
	defer listener.Close()
	log.Printf("Server is listening on port %s...\n", os.Args[2])
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go Handle(DB, conn)
	}
}
