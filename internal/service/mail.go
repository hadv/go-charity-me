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
	m := mail.NewV3Mail()
	address := "noreply@charity.me"
	name := "Charity Me"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)

	m.SetTemplateID(template)
	p := mail.NewPersonalization()
	user := data[0].(*model.User)
	token := data[1].(string)
	tos := []*mail.Email{
		mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email),
	}
	p.AddTos(tos...)
	p.SetDynamicTemplateData("reset_password_url", fmt.Sprintf("http://localhost:8081/reset-password?token=%s", token))
	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}
