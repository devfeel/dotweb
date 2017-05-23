package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/logger"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认启用，且默认为当前目录
	//app.SetLogger(NewYLog())
	app.SetEnabledLog(true)
	app.SetLogPath("d:/gotmp/xlog/xlog1/xlog2/")

	fmt.Println(logger.Logger())

	//开启development模式
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	logger.Logger().Debug("debug", "x")
	logger.Logger().Info("info", "x")
	logger.Logger().Warn("warn", "x")
	logger.Logger().Error("error", "x")
	_, err := ctx.WriteStringC(201, "index => ", ctx.RouterParams())
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
}

type chanLog struct {
	Content   string
	LogTarget string
}

type yLog struct {
	logRootPath    string
	logChan_Custom chan chanLog
	enabledLog     bool
}

//create new yLog
func NewYLog() *yLog {
	l := &yLog{logChan_Custom: make(chan chanLog, 10000)}
	go l.handleCustom()
	return l
}

const (
	defaultDateFormatForFileName = "2006_01_02"
	defaultDateLayout            = "2006-01-02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.999999"
	defaultTimeLayout            = "2006-01-02 15:04:05"
)

func (l *yLog) Debug(log string, logTarget string) {
	l.Log(log, logTarget, "debug")
}

func (l *yLog) Info(log string, logTarget string) {
	l.Log(log, logTarget, "info")
}

func (l *yLog) Warn(log string, logTarget string) {
	l.Log(log, logTarget, "warn")
}

func (l *yLog) Error(log string, logTarget string) {
	l.Log(log, logTarget, "error")
}

func (l *yLog) Log(log string, logTarget string, logLevel string) {
	if l.enabledLog {
		chanLog := chanLog{
			LogTarget: "yLog_" + logTarget + "_" + logLevel,
			Content:   log,
		}
		l.logChan_Custom <- chanLog
	}
}

//set log path
func (l *yLog) SetLogPath(rootPath string) {
	//设置日志根目录
	l.logRootPath = rootPath
	if !strings.HasSuffix(l.logRootPath, "/") {
		l.logRootPath = l.logRootPath + "/"
	}
}

//set enabled log
func (l *yLog) SetEnabledLog(enabledLog bool) {
	l.enabledLog = enabledLog
}

//处理日志内部函数
func (l *yLog) handleCustom() {
	for {
		log := <-l.logChan_Custom
		l.writeLog(log, "custom")
	}
}

func (l *yLog) writeLog(chanLog chanLog, level string) {
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
