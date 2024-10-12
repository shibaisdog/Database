package router

import (
	"log"
	"net"
)

var Con = "[ <- ]"
var Dis = "[ -> ]"

func LoggerConnet(request net.Conn) {
	log.Printf("%s <%s>", Con, request.RemoteAddr())
}

func LoggerDisconnet(request net.Conn) {
	log.Printf("%s <%s>", Dis, request.RemoteAddr())
	request.Close()
}

func LoggerDisconnetWithMessage(request net.Conn, message string) {
	log.Printf("%s <%s> : %s", Dis, request.RemoteAddr(), message)
	request.Close()
}
