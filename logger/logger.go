package logger

import (
	"github.com/devfeel/dotweb/framework/file"
)

const (
	LogLevel_Debug = "debug"
	LogLevel_Info  = "info"
	LogLevel_Warn  = "warn"
	LogLevel_Error = "error"
)

type AppLog interface {
	SetLogPath(logPath string)
	SetEnabledLog(enabledLog bool)
	Debug(log string, logTarget string)
	Info(log string, logTarget string)
	Warn(log string, logTarget string)
	Error(log string, logTarget string)
	Log(log string, logTarget string, logLevel string)
}

var (
	appLog         AppLog
	DefaultLogPath string
)

func Logger() AppLog {
	return appLog
}

func init() {
	DefaultLogPath = file.GetCurrentDirectory()
	appLog = NewXLog(DefaultLogPath)
	appLog.SetEnabledLog(true) //default enabled log
}
