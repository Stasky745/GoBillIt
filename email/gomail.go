package gomail

import "github.com/wneessen/go-mail"

type Email struct {
	From         string
	To           string
	Subject      string
	Body         string
	SmtpServer   string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	Attachment   string
}

func SendEmail(email Email) error {
	m := mail.NewMsg()
	if err := m.From(email.From); err != nil {
		return err
	}
	if err := m.To(email.To); err != nil {
		return err
	}
	m.AttachFile(email.Attachment)
	m.Subject(email.Subject)
	m.SetBodyString(mail.TypeTextPlain, email.Body)
	//
	// Deliver the mails via SMTP
	c, err := mail.NewClient(
		email.SmtpServer,
		mail.WithPort(email.SmtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(email.SmtpUsername),
		mail.WithPassword(email.SmtpPassword),
	)
	if err != nil {
		return err
	}
	err = c.DialAndSend(m)
	if err != nil {
		return err
	}
	c.Close()
	//
	return err
}
