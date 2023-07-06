package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)


func Error(msg ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	errorMsg := fmt.Sprintf("Error at %s:%d - ", file, line)
	errorMsg += fmt.Sprint(msg...)
	logrus.Error(errorMsg)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
