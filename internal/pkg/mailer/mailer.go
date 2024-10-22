package mailer

type Mailer interface {
	Host() string
	Set(to, subject, body string)
	Send() error
}
