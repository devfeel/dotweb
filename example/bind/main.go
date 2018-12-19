package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/framework/reflects"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	//这里仅为示例，默认情况下，开启的模式就是development模式
	app.SetDevelopmentMode()

	//使用json标签
	app.HttpServer.SetEnabledBindUseJsonTag(true)
	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置自定义绑定器
	app.HttpServer.SetBinder(newUserBinder())

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func TestBind(ctx dotweb.Context) error {
	type UserInfo struct {
		UserName string
		Sex      int
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.Bind(user); err != nil {
		errstr = err.Error()
	} else {

	}

	return ctx.WriteString("TestBind [" + errstr + "] " + fmt.Sprint(user))
}

func GetBind(ctx dotweb.Context) error {
	//type UserInfo struct {
	//	UserName string `form:"user"`
	//	Sex      int    `form:"sex"`
	//}
	type UserInfo struct {
		UserName string `json:"user"`
		Sex      int    `json:"sex"`
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.Bind(user); err != nil {
		errstr = err.Error()
	} else {

	}

	return ctx.WriteString("GetBind [" + errstr + "] " + fmt.Sprint(user))
}

func PostJsonBind(ctx dotweb.Context) error {
	type UserInfo struct {
		UserName string `json:"user"`
		Sex      int    `json:"sex"`
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.BindJsonBody(user); err != nil {
		errstr = err.Error()
	} else {

	}

	return ctx.WriteString("PostBind [" + errstr + "] " + fmt.Sprint(user))
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().POST("/", TestBind)
	server.Router().GET("/getbind", GetBind)
	server.Router().POST("/jsonbind", PostJsonBind)
}

type userBinder struct {
}

//Bind decode req.Body or form-value to struct
func (b *userBinder) Bind(i interface{}, ctx dotweb.Context) (err error) {
	fmt.Println("UserBind.Bind")
	req := ctx.Request()
	ctype := req.Header.Get(dotweb.HeaderContentType)
	if req.Body == nil {
		err = errors.New("request body can't be empty")
		return err
	}
	err = errors.New("request unsupported MediaType -> " + ctype)
	switch {
	case strings.HasPrefix(ctype, dotweb.MIMEApplicationJSON):
		err = json.Unmarshal(ctx.Request().PostBody(), i)
	case strings.HasPrefix(ctype, dotweb.MIMEApplicationXML):
		err = xml.Unmarshal(ctx.Request().PostBody(), i)
		//case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm),
		//	strings.HasPrefix(ctype, MIMETextHTML):
		//	err = reflects.ConvertMapToStruct(defaultTagName, i, ctx.FormValues())
	default:
		//check is use json tag, fixed for issue #91
		tagName := "form"
		if ctx.HttpServer().ServerConfig().EnabledBindUseJsonTag {
			tagName = "json"
		}
		//no check content type for fixed issue #6
		err = reflects.ConvertMapToStruct(tagName, i, ctx.Request().FormValues())
	}
	return err
}

//BindJsonBody default use json decode req.Body to struct
func (b *userBinder) BindJsonBody(i interface{}, ctx dotweb.Context) (err error) {
	fmt.Println("UserBind.BindJsonBody")
	if ctx.Request().PostBody() == nil {
		err = errors.New("request body can't be empty")
		return err
	}
	err = json.Unmarshal(ctx.Request().PostBody(), i)
	return err
}

func newUserBinder() *userBinder {
	return &userBinder{}
}
