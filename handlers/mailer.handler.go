package handlers

import (
	"log"

	"github.com/berezovskyi/domain-monitor/service"
	"github.com/labstack/echo/v4"
)

type MailerHandler struct {
	MailerService *service.MailerService
	Recipient     string
}

func NewMailerHandler(ms *service.MailerService, recipient string) *MailerHandler {
	// confirm that the mailer service is not nil
	if ms == nil {
		log.Fatal("🚨 Mailer service not properly initialized.")
	}

	return &MailerHandler{
		MailerService: ms,
		Recipient:     recipient,
	}
}

func (mh MailerHandler) HandleTestMail(c echo.Context) error {
	err := mh.MailerService.TestMail(mh.Recipient)
	if err != nil {
		log.Printf("❌ Failed to send test mail to %s: %s", mh.Recipient, err)
		return err
	}
	log.Printf("✅ Test mail sent successfully to %s", mh.Recipient)
	return c.JSON(200, "Mail sent")
}
