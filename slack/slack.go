package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// sample use
// 	slack.SendNotification(slack.Notification{
//	Title: "Error HTTP Request", URL: url, Headers: headers, Body: StreamToString(body),
//	Response: resp, ResponseCode: statusCode, Ctx: "GetHTTPRequestJSON", Error: err})

// Attachment model
type Attachment struct {
	Attachments []Payload `json:"attachments"`
	Channel     string    `json:"channel"`
}

// Payload model
type Payload struct {
	Text   string  `json:"text"`
	Color  string  `json:"color"`
	Fields []Field `json:"fields"`
}

// Field model
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Notification ...
type Notification struct {
	Title        string
	Headers      []map[string]string
	Request      string
	Body         string
	Ctx          string
	URL          string
	ResponseCode string
	Response     []byte
	Error        error
	Channel      string
	Username     string
	Indicator    string // info, warning, success
}

const (
	successColor        = "#36a64f"
	errorColor          = "#f44b42"
	warningColor        = "#f0ad4e"
	infoColor           = "#5bc0de"
	GenesisDevNotif     = "genesis-dev-notif"
	MigrationLoginext   = "C02GPREQ3CY"
	CheckTariffPrdNotif = "C02BDB5DP0Q"
	sauronDataNotif     = "C02MQ100DJR"
	InfoIndicator       = `info`
	WarningIndicator    = `warning`
	SuccessIndicator    = `success`
)

var (
	allowChannelPrd = map[string]bool{
		CheckTariffPrdNotif: true,
		MigrationLoginext:   true,
		sauronDataNotif:     true,
	}
)

func getCaller() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	var source string
	switch {
	case name != "":
		source = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		source = fmt.Sprintf("%v:%v", file, line)
	default:
		source = fmt.Sprintf("pc:%x", pc)
	}
	return source
}

// SendNotification to slack channel
func SendNotification(notification Notification) {
	isActive, _ := strconv.ParseBool(os.Getenv("SLACK_NOTIFIER"))
	if !isActive {
		return
	}

	stackTrace := getCaller()
	go func() {
		defer func() {
			if r := recover(); r != nil {

			}
		}()

		message := fmt.Sprintf("*%s*\n\n", notification.Title)
		if notification.URL != "" {
			message += fmt.Sprintf("*URL*\n\n```%s```\n\n", notification.URL)
		}
		if len(notification.Headers) > 0 {
			message += fmt.Sprintf("*Headers*\n\n```%s```\n\n", headersToString(notification.Headers))
		}
		if notification.Request != "" {
			message += fmt.Sprintf("*Request*\n\n```%s```\n\n", notification.Request)
		}
		if notification.Body != "" {
			message += fmt.Sprintf("*Body*\n\n```%s```\n\n", notification.Body)
		}
		if notification.ResponseCode != "" {
			message += fmt.Sprintf("*Response Code*\n\n```%s```\n\n", notification.ResponseCode)
		}
		if string(notification.Response) != "" {
			message += fmt.Sprintf("*Response*\n\n```%s```\n\n", notification.Response)
		}

		var slackPayload Payload
		slackPayload.Text = message
		/**
		 * Coloring
		 */
		switch notification.Indicator {
		case InfoIndicator:
			slackPayload.Color = infoColor
		case WarningIndicator:
			slackPayload.Color = warningColor
		case SuccessIndicator:
			slackPayload.Color = successColor
		}

		if notification.Request != "" {
			slackPayload.Fields = append(slackPayload.Fields, Field{
				Title: "Request",
				Value: fmt.Sprintf("`%s`", notification.Request),
				Short: true,
			})
		}

		if notification.Error != nil {
			slackPayload.Color = errorColor
			slackPayload.Text = fmt.Sprintf("%s\n*Error*: ```%s```", message, notification.Error.Error())
		}

		hostName, _ := os.Hostname()
		now := time.Now().Format(time.RFC3339)
		slackPayload.Fields = []Field{
			{
				Title: "Server",
				Value: hostName,
				Short: true,
			},
			{
				Title: "Environment",
				Value: os.Getenv("ENVIRONMENT"),
				Short: true,
			},
			{
				Title: "Context",
				Value: notification.Ctx,
				Short: true,
			},
			{
				Title: "Time",
				Value: now,
				Short: true,
			},
		}

		if notification.Error != nil {
			slackPayload.Fields = append(slackPayload.Fields, Field{
				Title: "Error Line Stack",
				Value: fmt.Sprintf("`%s`", stackTrace),
				Short: true,
			})
		}

		var slackAttachment Attachment
		slackAttachment.Attachments = append(slackAttachment.Attachments, slackPayload)

		/**
		 * Validation channel
		 */

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(true)
		encoder.Encode(slackAttachment)

		url := os.Getenv("WEBHOOK_SLACK")
		req, _ := http.NewRequest("POST", url, buffer)
		defer req.Body.Close()

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
	}()
}

func headersToString(headers []map[string]string) string {
	headerStr := ""
	for _, header := range headers {
		for key, value := range header {
			headerStr += fmt.Sprintf("%s=\"%s\"\n", key, value)
		}
	}

	return headerStr
}
