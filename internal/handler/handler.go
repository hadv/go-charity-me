package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// Response general json api response
type Response struct {
	Error *ErrorResponse `json:"error,omitempty"`
	Data  interface{}    `json:"data,omitempty"`
}

// ErrorResponse error json api response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// Render render request and return corresponding response
func (r *Response) Render(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func responseError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.Status(r, code)
	err := render.Render(w, r, &Response{
		Error: &ErrorResponse{
			Code:    code,
			Message: msg,
		},
	})
	if err != nil {
		log.Printf("cannot render response: %+v", err)
	}
}

func responseData(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := render.Render(w, r, &Response{
		Data: data,
	})
	if err != nil {
		log.Printf("error while responding data: %v\n%+v\n", err, data)
	}
}
