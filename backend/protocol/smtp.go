package protocol

import (
	"net/smtp"
)

type Smtp struct {
	host, port string
	email, pw  string
}

func NewSmtp(host, port, email, password string) *Smtp {
	return &Smtp{
		host: host, port: port, email: email, pw: password,
	}
}

func (s Smtp) Send(to []string, title, content string) error {
	message := []byte("Subject: " + title + "\r\n" + "\r\n" + content + "\r\n")
	auth := smtp.PlainAuth("", s.email, s.pw, s.host)
	err := smtp.SendMail(s.host+":"+s.port, auth, s.email, to, message)
	if err != nil {
		return err
	}
	return nil
}
