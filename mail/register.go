package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"api/logger"
	"config"
	"io"
	"net/smtp"
	"strings"
)

// Register sends a registration confirmation email to the specified address.
func Register(address string, token string) {
	c, err := smtp.Dial(config.Get().SMTP.Host + ":587")
	if err != nil {
		logger.Errorw("[SMTP] Failed to dial smtp host.", logger.Err(err))
		return
	}

	if err = c.Hello(config.Get().SMTP.Host); err != nil {
		logger.Errorw("[SMTP] Failed to send hello.", logger.Err(err))
		return
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: config.Get().SMTP.Host}
		if err = c.StartTLS(tlsConfig); err != nil {
			logger.Errorw("[SMTP] Failed to configure TLS connection.", logger.Err(err))
			return
		}
	}

	err = c.Auth(LoginAuth(config.Get().SMTP.Username, config.Get().SMTP.Password))
	if err != nil {
		logger.Errorw("[SMTP] Failed to authenticate with remote server.", logger.Err(err))
		return
	}

	if err = c.Mail(config.Get().SMTP.From); err != nil {
		logger.Errorw("[SMTP] Failed to set MAIL FROM.", logger.Err(err))
		return
	}

	if err = c.Rcpt(address); err != nil {
		logger.Errorw("[SMTP] Failed to set RCPT TO.", logger.Err(err))
		return
	}

	wc, err := c.Data()
	if err != nil {
		logger.Errorw("[SMTP] Failed to open data stream.", logger.Err(err))
		return
	}

	if _, err := io.Copy(wc, bytes.NewReader([]byte(fmt.Sprintf(`MIME-Version: 1.0
From: %s <%s>
To: <%s>
Subject: %s
Content-Type: text/plain; charset="utf-8"

Welcome to Ikuta!

Here is your registration confirmation link: https://egirls.me/user/register?token=%s`, config.Get().SMTP.Register.From, config.Get().SMTP.From, address, config.Get().SMTP.Register.Subject, token)))); err != nil {
		logger.Errorw("[SMTP] Failed to input data into data stream.", logger.Err(err))
		return
	}

	if err := wc.Close(); err != nil {
		if strings.Index(err.Error(), "200 ") != 0 {
			logger.Errorw("[SMTP] Received non-200 response when closing the data stream.", logger.Err(err))
			return
		}
	}

	if err := c.Quit(); err != nil {
		logger.Errorw("[SMTP] Failed to quit connection.", logger.Err(err))
		return
	}
}

type loginAuth struct {
	username, password string
}

// LoginAuth .
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start .
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Next .
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown fromServer")
		}
	}
	return nil, nil
}
