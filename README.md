# dotweb
基于go语言开发的web framework

启动代码：
    
    func StartServer() error {
	//初始化DotServer
	dotweb := dotweb.New()

	//设置dotserver日志目录
	dotweb.SetLogPath("/home/logs/wwwroot/")

	//设置路由
	InitRoute(dotweb)

	//启动监控服务
	pprofport := config.CurrentConfig.HttpServer.PProfPort
	go dotweb.StartPProfServer(pprofport)

	// 开始服务
	port := config.CurrentConfig.HttpServer.HttpPort
	innerLogger.Debug("dotweb.StartServer => " + strconv.Itoa(port))
	err := dotweb.StartServer(port)
	return err
    }

