package exception

import (
	"fmt"
	"runtime"
	"strings"
)

type Response struct {
	Status      bool        `json:"status"`
	Message     Message     `json:"message"`
	RealMessage string      `json:"real_message"`
	Data        interface{} `json:"data,omitempty"`
	Location    string      `json:"location,omitempty"`
}

type Message struct {
	Id string `json:"id"`
	En string `json:"en"`
}

func Error(err error, message Message, env string) *Response {
	pc, fn, line, _ := runtime.Caller(1)
	fnSplit := strings.Split(fn, "/")
	if env == "release" || env == "production" {
		return &Response{
			Status:  false,
			Message: message,
		}
	}
	return &Response{
		Status:      false,
		Message:     message,
		RealMessage: err.Error(),
		Location:    fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fnSplit[len(fnSplit)-1], line),
	}
}

func RouteNotFound() *Response {
	return &Response{
		Status: false,
		Message: Message{
			Id: "Halaman tidak ditemukan",
			En: "Page not found",
		},
	}
}

func Success(message Message, data interface{}) *Response {
	return &Response{
		Status:  true,
		Message: message,
		Data:    data,
	}
}
