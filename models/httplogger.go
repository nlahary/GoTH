package models

// HttpLogger is a logger that inherits from DefaultLogger
// used to log HTTP requests and responses
type HttpLogger struct {
	*DefaultLogger
}

const HttpLogSchema = `{
	"type": "record",
	"name": "HttpLog",
	"fields": [
		{"name": "method", "type": "string"},
		{"name": "status_code", "type": "int"},
		{"name": "url", "type": "string"},
		{"name": "body", "type": "string"}
	]
}`

type LogHTTP struct {
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	URL        string `json:"url"`
	Body       string `json:"body"`
}
