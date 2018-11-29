package adapter

// Context is ...
type Context interface {
	Param(string) string
	GetRawData() ([]byte, error)
	Status(int)
	Header(string, string)
	JSON(int, interface{})
}
