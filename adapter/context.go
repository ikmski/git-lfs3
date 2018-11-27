package adapter

// Context is ...
type Context interface {
	Param(string) string
	GetRawData() ([]byte, error)
	Status(int)
	Header(key, value, string)
	JSON(int, interface{})
}
