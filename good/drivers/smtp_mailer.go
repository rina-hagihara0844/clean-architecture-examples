package drivers

type SMTPMailer struct{}

func (m SMTPMailer) NotifyManagerNewRequest(id string) error { /* 実送信 */ return nil }
