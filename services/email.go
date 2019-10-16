package services

import (
	"edc-security-app/types"
	"errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

func SendEmail(req *types.MailRequest) error {
	log.Println("sending mail to ", req.Email)
	from := mail.NewEmail("Fast", "hi@fast.co")
	subject := req.Title
	to := mail.NewEmail(req.User, req.Email)
	message := mail.NewSingleEmail(from, subject, to, req.Body, req.Body)
	client := sendgrid.NewSendClient(os.Getenv("SG_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Printf("email sent to %s; Response => %v\n", req.Email, response)
	if response.StatusCode > 299 {
		return errors.New(response.Body)
	}

	return nil
}