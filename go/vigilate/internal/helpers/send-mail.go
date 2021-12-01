package helpers

import "github.com/tsawler/vigilate/internal/channeldata"

// SendEmail sends an email
func SendEmail(mailMessage channeldata.MailData) {
	job := channeldata.MailJob{MailMessage: mailMessage}
	app.MailQueue <- job
}
