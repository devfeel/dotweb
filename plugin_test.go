package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/test"
	"testing"
)

type testPlugin struct {
}

func (p *testPlugin) Name() string {
	return "test"
}
func (p *testPlugin) Run() error {
	fmt.Println(p.Name(), "runing")
	//panic("error test run")
	return nil
}
func (p *testPlugin) IsValidate() bool {
	return true
}

func TestNotifyPlugin_Name(t *testing.T) {
	app := newConfigDotWeb()
	//fmt.Println(app.Config.ConfigFilePath)
	p := NewDefaultNotifyPlugin(app)
	needShow := "NotifyPlugin"
	test.Equal(t, needShow, p.Name())
}

func TestNotifyPlugin_IsValidate(t *testing.T) {
	app := newConfigDotWeb()
	p := NewDefaultNotifyPlugin(app)
	needShow := true
	test.Equal(t, needShow, p.IsValidate())
}
