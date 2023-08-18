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

var CreateType = "create"
var DeleteteType = "delete"
var UpdateType = "update"
var GetType = "get"

func messageAction(sub string) Message {
	return Message{
		Id: fmt.Sprintf("gagal %s data!", sub),
		En: fmt.Sprintf("failed %s data!", sub),
	}
}

func messageActionFind() Message {
	return Message{
		Id: "data tidak ditemukan",
		En: "data not found",
	}
}

func Error(err error, typeAction string, env string) *Response {
	var message Message
	switch typeAction {
	case CreateType:
		message = messageAction(CreateType)
	case DeleteteType:
		message = messageAction(DeleteteType)
	case UpdateType:
		message = messageAction(UpdateType)
	case GetType:
		message = messageActionFind()
	default:
		message = Message{}
	}

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
