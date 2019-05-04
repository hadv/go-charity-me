package service

import (
	"fmt"

	"github.com/hadv/go-charity-me/internal/model"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"
)

func sendMail(template string, data []interface{}, dynamicTemplate func(string, []interface{}) []byte) {
	request := sendgrid.GetRequest(viper.GetString("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = dynamicTemplate(template, data)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func createForgotPasswordEmailFromTemplate(template string, data []interface{}) []byte {
	email := mail.NewV3Mail()
	from := mail.NewEmail("Charity Me", "noreply@charity.me")
	email.SetFrom(from)

	email.SetTemplateID(template)
	persona := mail.NewPersonalization()
	user := data[0].(*model.User)
	token := data[1].(string)
	tos := []*mail.Email{
		mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email),
	}
	persona.AddTos(tos...)
	persona.SetDynamicTemplateData("reset_password_url", fmt.Sprintf("http://localhost:8081/reset-password?token=%s", token))
	email.AddPersonalizations(persona)
	return mail.GetRequestBody(email)
}

func createVerifyEmailFromTemplate(template string, data []interface{}) []byte {
	email := mail.NewV3Mail()
	from := mail.NewEmail("Charity Me", "noreply@charity.me")
	email.SetFrom(from)

	email.SetTemplateID(template)
	persona := mail.NewPersonalization()
	user := data[0].(*model.User)
	token := data[1].(string)
	tos := []*mail.Email{
		mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email),
	}
	persona.AddTos(tos...)
	persona.SetDynamicTemplateData("verify_email_url", fmt.Sprintf("http://localhost:8081/verify-email?token=%s", token))
	email.AddPersonalizations(persona)
	return mail.GetRequestBody(email)
}
