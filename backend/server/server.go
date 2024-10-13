package main

import (
	"cloud-logic/database/memory"
	"cloud-logic/protocol"
	"cloud-logic/server/logger"
	"cloud-logic/server/router"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func Handle(smtp *protocol.Smtp, db *memory.Storage, conn net.Conn) {
	defer conn.Close()
	client := router.Client{
		Login: false,
	}
	logger.LoggerConnet(conn)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				logger.LoggerDisconnetWithMessage(conn, "Timeout")
				break
			} else if err.Error() == "EOF" {
				logger.LoggerDisconnet(conn)
				break
			}
			logger.LoggerDisconnetWithMessage(conn, err.Error())
			break
		}
		go router.Router(smtp, db, &client, conn, buffer[:n])
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	SMTP := protocol.NewSmtp(os.Getenv("Host"), os.Getenv("Port"), os.Getenv("Email"), os.Getenv("Password"))
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
		go Handle(SMTP, DB, conn)
	}
}
