package dotweb

import (
	"fmt"
	"testing"

	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/test"
)

// 以下为功能测试

// 测试RunMode函数无配置文件时的返回值
func Test_RunMode_1(t *testing.T) {
	app := New()
	runMode := app.RunMode()
	t.Log("RunMode:", runMode)
}

// 测试RunMode函数有配置文件时的返回值
func Test_RunMode_2(t *testing.T) {
	runModes := []string{"dev", "development", "prod", "production"}

	app := New()
	for _, value := range runModes {
		app.Config.App.RunMode = value
		runMode := app.RunMode()
		t.Log("runModes value:", value, "RunMode:", runMode)
	}
}

//测试IsDevelopmentMode函数
func Test_IsDevelopmentMode_1(t *testing.T) {
	app := New()
	app.Config.App.RunMode = "development"
	b := app.IsDevelopmentMode()
	test.Equal(t, true, b)
	t.Log("Run IsDevelopmentMode :", b)
}

func Test_IsDevelopmentMode_2(t *testing.T) {
	app := New()
	app.Config.App.RunMode = "production"
	b := app.IsDevelopmentMode()
	t.Log("Run IsDevelopmentMode :", b)
}

func TestDotWeb_UsePlugin(t *testing.T) {
	app := newConfigDotWeb()
	app.UsePlugin(new(testPlugin))
	app.UsePlugin(NewDefaultNotifyPlugin(app))
	fmt.Println(app.pluginMap)
	app.StartServer(8081)
}

func newConfigDotWeb() *DotWeb {
	app := New()
	appConfig, err := config.InitConfig("config/testdata/dotweb.conf", "xml")
	if err != nil {
		fmt.Println("dotweb.InitConfig error => " + fmt.Sprint(err))
		return nil
	}
	app.Logger().SetEnabledConsole(true)
	app.SetConfig(appConfig)
	return app
}
