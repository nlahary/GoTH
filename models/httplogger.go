package models

import (
	kafka "github.com/nlahary/website/kafka"
)

type HttpLogger struct {
	kafka.Producer
}

const HttpLogSchema = `{
	"type": "record",
	"name": "HttpLog",
	"fields": [
		{"name": "method", "type": "string"},
		{"name": "status_code", "type": "int"},
		{"name": "url", "type": "string"},
		{"name": "body", "type": "string"},
		{"name": "response_time", "type": "int"}
	]
}`

type LogHTTP struct {
	Method       string `json:"method"`
	StatusCode   int    `json:"status_code"`
	URL          string `json:"url"`
	Body         string `json:"body"`
	ResponseTime int64  `json:"response_time"`
}
