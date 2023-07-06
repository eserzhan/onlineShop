package logger

type Logger interface {
	Error(msg string, params map[string]interface{})
}
