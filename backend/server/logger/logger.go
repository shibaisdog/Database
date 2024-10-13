package logger

import (
	"log"
	"net"
)

var Con = Fg_Green + "[ <- ]" + Fg_BrightYellow
var Dis = Fg_Red + "[ -> ]" + Fg_BrightYellow

func LoggerConnet(request net.Conn) {
	log.Printf("%s <%s> %s", Con, request.RemoteAddr(), Reset)
}

func LoggerDisconnet(request net.Conn) {
	log.Printf("%s <%s> %s", Dis, request.RemoteAddr(), Reset)
	request.Close()
}

func LoggerDisconnetWithMessage(request net.Conn, message string) {
	log.Printf("%s <%s> : %s %s", Dis, request.RemoteAddr(), message, Reset)
	request.Close()
}
