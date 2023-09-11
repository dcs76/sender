package sender

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

const (
	// SMTPServer for gmail
	SMTPServer = "smtp.gmail.com"
)

// Sender instances can send email using smtp.gmail.com
type Sender struct {
	User     string
	Password string
}

// New returns a Sender
func New(Username, Password string) Sender {
	return Sender{Username, Password}
}

// sendMail sends either a plain or an html formatted body message
func (s Sender) sendMail(Dest []string, Subject, bodyMessage string) error {
	msg := "From: " + s.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", s.User, s.Password, SMTPServer),
		s.User, Dest, []byte(msg))

	return err
}

// writeEmail formats the body message according to contentType
func (s Sender) writeEmail(dest []string, contentType, subject, bodyMessage string) string {
	header := make(map[string]string)
	header["From"] = s.User

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

// writeHTMLEmail formats a message with content type "text/html"
func (s Sender) writeHTMLEmail(dest []string, subject, bodyMessage string) string {
	return s.writeEmail(dest, "text/html", subject, bodyMessage)
}

// writePlainEmail formats a message with content type "text/plain"
func (s Sender) writePlainEmail(dest []string, subject, bodyMessage string) string {
	return s.writeEmail(dest, "text/plain", subject, bodyMessage)
}

// SendPlainMessage sends a plain email message to the list in dest
func (s Sender) SendPlainMessage(dest []string, subject, bodyMessage string) error {
	msg := s.writePlainEmail(dest, subject, bodyMessage)
	return s.sendMail(dest, subject, msg)
}
