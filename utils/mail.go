package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/inchingforward/mnmnt/models"
	"github.com/mailgun/mailgun-go"
)

// SendEmail sends a "New Monument memory submitted" email along with an
// approval link.  The email is sent to the recipients defined in the
// MONUMENT_ADMIN_EMAIL environment variable.  If any mail-related
// environment variable is not set, no email is sent.
func SendEmail(memory models.Memory) {
	domain := os.Getenv("MONUMENT_MAILGUN_DOMAIN")
	prvKey := os.Getenv("MONUMENT_MAILGUN_PRIVATE_KEY")
	pubKey := os.Getenv("MONUMENT_MAILGUN_PUBLIC_KEY")
	address := os.Getenv("MONUMENT_ADMIN_EMAIL")
	mnmntHost := os.Getenv("MONUMENT_HOST")

	if domain == "" || prvKey == "" || pubKey == "" || address == "" || mnmntHost == "" {
		log.Println("Missing mail environment variables...not sending")
		return
	}

	approvalLink := fmt.Sprintf("%v/memories/approve/%v", mnmntHost, memory.ApprovalUUID)
	subject := "New Monument memory submitted"
	body := fmt.Sprintf("%v:\n\n%v\n\n-%v\n\nApproval link: %v", memory.Title, memory.Details, memory.Author, approvalLink)

	gun := mailgun.NewMailgun(domain, prvKey)
	msg := gun.NewMessage(address, subject, body, address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, id, err := gun.Send(ctx, msg)
	log.Printf("mailer response: %v, message: %v, error: %v\n", id, response, err)
}
