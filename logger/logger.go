package logger

import (
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/devfeel/dotweb/framework/file"
)

const (
	// LogLevelDebug raw log level
	LogLevelRaw = "RAW"
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
	IsEnabledLog() bool
	Print(log string, logTarget string)
	Raw(log string, logTarget string)
	Debug(log string, logTarget string)
	Info(log string, logTarget string)
	Warn(log string, logTarget string)
	Error(log string, logTarget string)
}

var (
	DefaultLogPath        string
	DefaultEnabledLog     bool = false
	DefaultEnabledConsole bool = false
)

func NewAppLog() AppLog {
	if DefaultLogPath == "" {
		DefaultLogPath = file.GetCurrentDirectory() + "/logs"
	}
	appLog := NewXLog()
	appLog.SetLogPath(DefaultLogPath)               // set default log path
	appLog.SetEnabledLog(DefaultEnabledLog)         // set default enabled log
	appLog.SetEnabledConsole(DefaultEnabledConsole) // set default enabled console output
	return appLog
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
