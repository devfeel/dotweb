package logger

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type chanLog struct {
	Content   string
	LogTarget string
}

type xLog struct {
	logRootPath    string
	logChan_Custom chan chanLog
	enabledLog     bool
}

//create new xLog
func NewXLog(logPath string) *xLog {
	l := &xLog{logChan_Custom: make(chan chanLog, 10000)}
	//设置日志根目录
	l.SetLogPath(logPath)
	go l.handleCustom()
	return l
}

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultDateLayout            = "2006-01-02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
	defaultTimeLayout            = "2006-01-02 15:04:05"
)

func (l *xLog) Debug(log string, logTarget string) {
	l.Log(log, logTarget, "debug")
}

func (l *xLog) Info(log string, logTarget string) {
	l.Log(log, logTarget, "info")
}

func (l *xLog) Warn(log string, logTarget string) {
	l.Log(log, logTarget, "warn")
}

func (l *xLog) Error(log string, logTarget string) {
	l.Log(log, logTarget, "error")
}

func (l *xLog) Log(log string, logTarget string, logLevel string) {
	if l.enabledLog {
		chanLog := chanLog{
			LogTarget: logTarget + "_" + logLevel,
			Content:   log,
		}
		l.logChan_Custom <- chanLog
	}
}

//set log path
func (l *xLog) SetLogPath(rootPath string) {
	//设置日志根目录
	l.logRootPath = rootPath
	if !strings.HasSuffix(l.logRootPath, "/") {
		l.logRootPath = l.logRootPath + "/"
	}
}

//set enabled log
func (l *xLog) SetEnabledLog(enabledLog bool) {
	l.enabledLog = enabledLog
}

//处理日志内部函数
func (l *xLog) handleCustom() {
	for {
		log := <-l.logChan_Custom
		l.writeLog(log, "custom")
	}
}

func (l *xLog) writeLog(chanLog chanLog, level string) {
	filePath := l.logRootPath + chanLog.LogTarget
	switch level {
	case "custom":
		filePath = filePath + "_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
		break
	}
	log := time.Now().Format(defaultFullTimeLayout) + " " + chanLog.Content
	writeFile(filePath, log)
}

func writeFile(logFile string, log string) {
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	logstr := log + "\r\n"
	file, err := os.OpenFile(logFile, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(logFile, err)
		return
	}
	//fmt.Print(logstr)
	file.WriteString(logstr)
}
