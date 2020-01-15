package adapter

import "io"

// Context is ...
type Context interface {
	GetHeader(string) string
	GetParam(string) string
	GetRawData() ([]byte, error)
	SetStatus(int)
	SetHeader(string, string)

	GetResponseWriter() io.Writer
	GetRequestReader() io.Reader
}
