package email

import (
	"github.com/Stasky745/go-libs/log"
	"github.com/wneessen/go-mail"
)

type Smtp struct {
	Server   string `json:"server,omitempty"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type Email struct {
	From       string   `json:"from,omitempty"`
	To         []string `json:"to,omitempty"`
	Cc         []string `json:"cc,omitempty"`
	Bcc        []string `json:"bcc,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	Body       string   `json:"body,omitempty"`
	Smtp       Smtp     `json:"smtp,omitempty"`
	Attachment string
}

func SendEmail(email Email) error {
	m := mail.NewMsg()

	if err := m.From(email.From); err != nil {
		return err
	}
	if err := m.To(email.To...); err != nil {
		return err
	}
	if err := m.Cc(email.Cc...); err != nil {
		return err
	}
	if err := m.Bcc(email.Bcc...); err != nil {
		return err
	}
	m.AttachFile(email.Attachment)
	m.Subject(email.Subject)
	m.SetBodyString(mail.TypeTextHTML, email.Body)

	// Deliver the mails via SMTP
	c, err := mail.NewClient(
		email.Smtp.Server,
		mail.WithPort(email.Smtp.Port),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(email.Smtp.Username),
		mail.WithPassword(email.Smtp.Password),
	)
	if log.CheckErr(err, false, "can't create new mail client", "email", email) {
		return err
	}
	err = c.DialAndSend(m)
	log.CheckErr(err, false, "can't send email")

	c.Close()
	return err
}
