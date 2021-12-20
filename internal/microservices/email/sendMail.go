package main

import (
	"github.com/sirupsen/logrus"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

const (
	encryptionType = mail.EncryptionSSL
	connectTimeout = 10 * time.Second
	sendTimeout    = 10 * time.Second
)

type HTMLString string

type SenderEmail struct {
	broadcast chan *message
	stop      chan bool
	client    *mail.SMTPServer
	log       *logrus.Entry
}

type message struct {
	to      []string
	message HTMLString
}

// Some variables to connect and the body.
var (
	htmlBody = `<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
		<title>Hello Gophers!</title>
	</head>
	<body>
		<p>This is the <b>Go gopher</b>.</p>
		<p><img src="cid:Gopher.png" alt="Go gopher" /></p>
		<p>Image created by Renee French</p>
	</body>
</html>`

	host     = "localhost"
	port     = 25
	username = "test@example.com"
	password = "santiago"
)

func NewEmailSender(host string, port int64, emailFrom string, password string, log *logrus.Entry) *SenderEmail {
	res := &SenderEmail{
		broadcast: make(chan *message),
		stop:      make(chan bool),
		client:    mail.NewSMTPClient(),
		log:       log,
	}
	res.client.Host = host
	res.client.Port = int(port)
	res.client.Username = emailFrom
	res.client.Password = password
	res.client.Encryption = encryptionType
	res.client.ConnectTimeout = connectTimeout
	res.client.SendTimeout = sendTimeout
	res.client.Authentication = mail.AuthPlain
	res.client.KeepAlive = true
	return res
}

func (h *SenderEmail) SendMessage(to []string, hsg HTMLString) {
	h.broadcast <- &(message{to: to, message: hsg})
}

func (h *SenderEmail) Stop() {
	h.stop <- true
}

func (h *SenderEmail) sendMessage(msg *message) {
	conn, err := h.client.Connect()
	if err != nil {
		h.log.Errorf("Some error when try connect to smtp client %s", err)
	}
	defer func(conn *mail.SMTPClient) {
		_ = conn.Quit()
		_ = conn.Close()
	}(conn)

	for _, to := range msg.to {
		err = h.sendEmail(string(msg.message), to, conn)
		if err != nil {
			h.log.Error("Expected nil, got", err, "sending email")
		}
	}
}

func (h *SenderEmail) Run() {
	for {
		select {
		case msg, ok := <-h.broadcast:
			if ok {
				h.sendMessage(msg)
			}
			break
		case <-h.stop:
			return
		default:
			break
		}
	}
}

func (h *SenderEmail) sendEmail(htmlBody string, to string, smtpClient *mail.SMTPClient) error {
	email := mail.NewMSG()

	email.SetFrom("noreply@pyaterocka-team.site").
		AddTo(to).
		SetSubject("Patreon email")

	email.GetFrom()
	email.SetBody(mail.TextHTML, htmlBody)

	email.SetPriority(mail.PriorityHigh)

	if email.Error != nil {
		return email.Error
	}

	return email.Send(smtpClient)
}