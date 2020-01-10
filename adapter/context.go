package adapter

// Context is ...
type Context interface {
	GetHeader(string) string
	GetParam(string) string
	GetRawData() ([]byte, error)
	SetStatus(int)
	SetHeader(string, string)
	SetJson(int, interface{})
}
