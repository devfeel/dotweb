package logger

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/devfeel/dotweb/framework/file"
)

const (
	// LogLevelDebug debug log level
	LogLevelDebug = "DEBUG"
	// LogLevelInfo info log level
	LogLevelInfo = "INFO"
	// LogLevelWarn warn log level
	LogLevelWarn = "WARN"
	// LogLevelError error log level
	LogLevelError = "ERROR"
)

type AppLog interface {
	SetLogPath(logPath string)
	SetEnabledConsole(enabled bool)
	SetEnabledLog(enabledLog bool)
	Debug(log string, logTarget string)
	Print(log string, logTarget string)
	Info(log string, logTarget string)
	Warn(log string, logTarget string)
	Error(log string, logTarget string)
}

var (
	appLog         AppLog
	DefaultLogPath string
	EnabledLog     bool = false
	EnabledConsole bool = false
)

func Logger() AppLog {
	return appLog
}

// SetLogPath set log path
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

// SetEnabledLog set enabled log
func SetEnabledLog(isLog bool) {
	EnabledLog = isLog
	if appLog != nil {
		appLog.SetEnabledLog(isLog)
	}
}

// SetEnabledConsole set enabled Console output
func SetEnabledConsole(enabled bool) {
	EnabledConsole = enabled
	if appLog != nil {
		appLog.SetEnabledConsole(enabled)
	}
}

func InitLog() {
	if DefaultLogPath == "" {
		DefaultLogPath = file.GetCurrentDirectory() + "/logs"
	}
	if appLog == nil {
		appLog = NewXLog()
	}

	SetLogPath(DefaultLogPath)        // set default log path
	SetEnabledLog(EnabledLog)         // set default enabled log
	SetEnabledConsole(EnabledConsole) // set default enabled console output
}

// Log content
// fileName source file name
// line line number in source file
// fullPath full path of source file
// funcName function name of caller
type logContext struct {
	fileName string
	line     int
	fullPath string
	funcName string
}

// priting
// skip=0  runtime.Caller
// skip=1  runtime/proc.c: runtime.main
// skip=2  runtime/proc.c: runtime.goexit
//
// Process startup procedure of a go program:
// 1.runtime.goexit is the actual entry point(NOT main.main)
// 2.then runtime.goexit calls runtime.main
// 3.finally runtime.main calls user defined main.main
func callerInfo(skip int) (ctx *logContext, err error) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("error  during runtime.Callers")
	}

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		return nil, errors.New("error during runtime.FuncForPC")
	}

	funcName := funcInfo.Name()
	if strings.HasPrefix(funcName, ".") {
		funcName = funcName[strings.Index(funcName, "."):]
	}

	ctx = &logContext{
		funcName: filepath.Base(funcName),
		line:     line,
		fullPath: file,
		fileName: filepath.Base(file),
	}

	return ctx, nil

}
