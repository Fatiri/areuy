package exception

import (
	"fmt"
	"runtime"
	"strings"
)

type Response struct {
	Status   bool        `json:"status"`
	Message  Message     `json:"message"`
	Error    interface{} `json:"error,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Location string      `json:"location,omitempty"`
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
	if err != nil {
		return &Response{
			Status:   false,
			Message:  message,
			Error:    err.Error(),
			Location: fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fnSplit[len(fnSplit)-1], line),
		}
	} else {
		return &Response{
			Status:   false,
			Message:  message,
			Location: fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fnSplit[len(fnSplit)-1], line),
		}
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
