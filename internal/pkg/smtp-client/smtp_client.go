package smtpclient

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
)

type SMTPClient struct {
	host      string
	port      string
	auth      smtp.Auth
	tlsConfig *tls.Config
	client    *smtp.Client
}

func New(
	host string,
	port string,
	username string,
	password string,
) (*SMTPClient, error) {
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)

	tlsConfig := tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	}

	return &SMTPClient{
		host:      host,
		port:      port,
		auth:      auth,
		tlsConfig: &tlsConfig,
	}, nil
}

func (s *SMTPClient) Connect() error {
	client, err := smtp.Dial(s.host + ":" + s.port)

	if err != nil {
		return err
	}

	s.client = client
	err = s.client.StartTLS(s.tlsConfig)

	if err != nil {
		return err
	}

	return s.client.Auth(s.auth)
}

func (s *SMTPClient) CreateMessage(
	from mail.Address,
	to mail.Address,
	title string,
	body string,
) []byte {
	headers := make(map[string]string)

	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = title
	headers["Content-Type"] = "text/html"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	return []byte(message)
}

func (s *SMTPClient) SendEmail(
	from mail.Address,
	to mail.Address,
	message []byte,
) error {
	err := s.client.Mail(from.Address)

	if err != nil {
		return err
	}

	err = s.client.Rcpt(to.Address)

	if err != nil {
		return err
	}

	writer, err := s.client.Data()

	if err != nil {
		return err
	}

	_, err = writer.Write(message)

	if err != nil {
		return err
	}

	return writer.Close()
}

func (s *SMTPClient) Disconnect() error {
	return s.client.Quit()
}
