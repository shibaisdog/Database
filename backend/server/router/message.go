package router

import (
	"cloud-logic/database"
	"net"
)

func SendErrorMessage(request net.Conn, state int, err string) {
	msg, _ := database.ByteJSON(&map[string]interface{}{
		"state": state,
		"error": err,
	})
	request.Write(msg)
}

func SendErrorMessageWithDisconnet(request net.Conn, state int, err string) {
	msg, _ := database.ByteJSON(&map[string]interface{}{
		"state": state,
		"error": err,
	})
	request.Write(msg)
	request.Close()
}

func SendMessage(request net.Conn, state int, message string) {
	msg, _ := database.ByteJSON(&map[string]interface{}{
		"state": state,
		"data":  message,
	})
	request.Write(msg)
}
