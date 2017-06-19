package dotweb

import (
	"testing"
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
	t.Log("Run IsDevelopmentMode :", b)
}

func Test_IsDevelopmentMode_2(t *testing.T) {
	app := New()
	app.Config.App.RunMode = "production"
	b := app.IsDevelopmentMode()
	t.Log("Run IsDevelopmentMode :", b)
}