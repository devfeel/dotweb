// logger
package logger

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type ChanLog struct {
	Content   string
	LogTarget string
}

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultDateLayout            = "2006-01-02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
	defaultTimeLayout            = "2006-01-02 15:04:05"
)

var (
	logChan_Custom chan ChanLog
)

var (
	logRootPath string
)

func init() {
	logChan_Custom = make(chan ChanLog, 10000)
}

func Debug(log string, logTarget string) {
	Log(log, logTarget, "debug")
}

func Info(log string, logTarget string) {
	Log(log, logTarget, "info")
}

func Warn(log string, logTarget string) {
	Log(log, logTarget, "warn")
}

func Error(log string, logTarget string) {
	Log(log, logTarget, "error")
}

func Log(log string, logTarget string, logLevel string) {
	chanLog := ChanLog{
		LogTarget: logTarget + "_" + logLevel,
		Content:   log,
	}
	logChan_Custom <- chanLog
}

//开启日志处理器
func StartLogHandler(rootPath string) {
	//设置日志根目录
	logRootPath = rootPath
	if !strings.HasSuffix(logRootPath, "/") {
		logRootPath = logRootPath + "/"
	}
	go handleCustom()
}

//处理日志内部函数
func handleCustom() {
	for {
		log := <-logChan_Custom
		writeLog(log, "custom")
	}
}

func writeLog(chanLog ChanLog, level string) {
	filePath := logRootPath + chanLog.LogTarget
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
	file.WriteString(logstr)
}
