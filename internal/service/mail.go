package service

import (
	"fmt"
	"log"
	"os"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendMail(user *model.User, subject, text, html string) {
	from := mail.NewEmail("Charity Me", "noreply@charity.me")
	to := mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email)
	message := mail.NewSingleEmail(from, subject, to, text, html)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
