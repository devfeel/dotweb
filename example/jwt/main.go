package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/middleware/jwt"
	"strconv"
	"time"
)

const JwtContextKey = "jwtuser"

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	user, exists := ctx.Items().Get(JwtContextKey)
	_, err := ctx.WriteString("custom jwt context => ", user, exists)
	return err
}

func Login(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	config := parseJwtConfig(ctx.AppContext().Get("CustomJwtConfig"))
	if config == nil {
		_, err := ctx.WriteString("custom login failed, token config not exists")
		return err
	}
	m := make(map[string]interface{})
	m["userid"] = "loginuser"
	m["userip"] = ctx.RemoteIP()
	token, err := jwt.GeneratorToken(config, m)
	if err != nil || token == "" {
		_, err := ctx.WriteString("custom login failed, token create failed, ", err.Error())
		return err
	}

	ctx.SetCookieValue(config.Name, token, 0)
	_, err = ctx.WriteString("custom login is ok, token => ", token)
	return err
}

func Logout(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	config := parseJwtConfig(ctx.AppContext().Get("CustomJwtConfig"))
	if config == nil {
		_, err := ctx.WriteString("logout failed, token config not exists")
		return err
	}
	ctx.RemoveCookie(config.Name)
	_, err := ctx.WriteString("logout is ok")
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index).Use(NewCustomJwt(server.DotApp))
	server.Router().GET("/Login", Login)
	server.Router().GET("/Logout", Logout)
}

func NewSimpleJwt(app *dotweb.DotWeb) dotweb.Middleware {
	option := &jwt.Config{
		SigningKey: []byte("devfeel/dotweb"), //must input
		//use cookie
		Extractor: jwt.ExtractorFromCookie,
	}
	app.AppContext.Set("SimpleJwtConfig", option)
	return jwt.NewJWT(option)
}

func NewCustomJwt(app *dotweb.DotWeb) dotweb.Middleware {
	option := &jwt.Config{
		TTL:           time.Minute * 10,         //default is 24 hour
		ContextKey:    JwtContextKey,            //default is dotuser
		SigningKey:    []byte("devfeel/dotweb"), //must input
		SigningMethod: jwt.SigningMethodHS256,   //default is SigningMethodHS256
		ExceptionHandler: func(ctx dotweb.Context, err error) {
			//TODO:log err info
			ctx.WriteString("no authorization, please login first")
		},
		AddonValidator: func(config *jwt.Config, ctx dotweb.Context) error {
			//example: check user ip
			user, exists := ctx.Items().Get(JwtContextKey)
			if !exists {
				return errors.New("no token exists")
			}
			fmt.Println(user)
			jwtUserIp := user.(map[string]interface{})["userip"].(string)
			requestIp := ctx.RemoteIP()
			fmt.Println("jwtUserIp", jwtUserIp, " requestIp:", requestIp)
			if jwtUserIp != requestIp {
				return errors.New("ip is not match")
			}
			return nil
		},
		//use cookie
		Extractor: jwt.ExtractorFromCookie,
	}

	app.AppContext.Set("CustomJwtConfig", option)

	return jwt.NewJWT(option)
}

func parseJwtConfig(c interface{}, exists bool) (config *jwt.Config) {
	if c == nil || !exists {
		return nil
	}
	config = c.(*jwt.Config)
	return config
}
