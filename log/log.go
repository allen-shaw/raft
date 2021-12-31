package log

import "log"

func Debug(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Info(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Warn(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Error(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}
