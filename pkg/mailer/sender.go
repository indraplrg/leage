package mailer

import (
	"fmt"
	"os"
	"share-notes-app/configs"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	from string
	baseUrl string
}

func NewMailer(config *configs.Config) *Mailer {
 d := gomail.NewDialer(
	config.SMTP.Host,
	config.SMTP.Port,
	os.Getenv("APP_SMTP_AUTH_EMAIL"),
	os.Getenv("APP_SMTP_PASSWORD_EMAIL"),
 )

 return &Mailer{
	dialer: d,
	from: config.SMTP.SenderName,
	baseUrl: config.BaseUrl,
 }
}

func (m *Mailer) SendVerification(to string, token string) error {
	msg := gomail.NewMessage()

	verifyUrl := fmt.Sprintf(
		"%s/api/auth/verify-email/%s",
		m.baseUrl,
		token,
	)

	body := fmt.Sprintf(`
		<html>
		<body>
			<p>Hi 👋</p>
			<p>Click the button below to verify your email:</p>

			<a href="%s"
			   style="
			     display:inline-block;
			     padding:12px 20px;
			     background-color:#4f46e5;
			     color:#ffffff;
			     text-decoration:none;
			     border-radius:6px;
			     font-weight:600;
			   ">
			   Verify Email
			</a>

			<p style="margin-top:20px;font-size:12px;color:#666;">
			  This link will expire in 24 hours.
			</p>
		</body>
		</html>
	`, verifyUrl)

	// set mailer
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "hei kontol")
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}