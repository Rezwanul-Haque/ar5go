package domain

type ILogger interface {
	Debug(msg string, data interface{})
	Error(msg string, err error)
	Info(msg string)
	Warn(msg string)
	Fatal(msg string)
	Panic(msg string)
}
