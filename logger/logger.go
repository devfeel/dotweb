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
	EnabledLog     bool = true
)

func Logger() AppLog {
	return appLog
}

func SetLogger(logger AppLog) {
	appLog = logger
	logger.SetLogPath(DefaultLogPath)
	logger.SetEnabledLog(EnabledLog)
}

func SetLogPath(path string) {
	DefaultLogPath = path
	if appLog != nil {
		appLog.SetLogPath(path)
	}
}

func SetEnabledLog(isLog bool) {
	EnabledLog = isLog
	if appLog != nil {
		appLog.SetEnabledLog(isLog)
	}
}

func InitLog() {
	if DefaultLogPath == "" {
		DefaultLogPath = file.GetCurrentDirectory()
	}
	if appLog == nil {
		appLog = NewXLog()
	}

	SetLogPath(DefaultLogPath) //set default log path
	SetEnabledLog(EnabledLog)  //set default enabled log
}
