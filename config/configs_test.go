package config


//运行以下用例需要在edit configuration中将working dir改成dotweb目录下,不能在当前目录
import (
	"testing"
	"github.com/devfeel/dotweb/test"
)


func TestInitConfig(t *testing.T) {
	conf,err:=InitConfig("example/config/dotweb.json.conf","json")

	test.Nil(t,err)
	test.NotNil(t,conf)
	test.NotNil(t,conf.App)
	test.NotNil(t,conf.App.LogPath)
	test.NotNil(t,conf.ConfigSet)
	test.Equal(t,4,	conf.ConfigSet.Len())
}

//该测试方法报错...
//是xml问题还是代码问题?
func TestInitConfigWithXml(t *testing.T) {
	conf,err:=InitConfig("example/config/dotweb.conf","xml")

	test.Nil(t,err)
	test.NotNil(t,conf)
	test.NotNil(t,conf.App)
	test.NotNil(t,conf.App.LogPath)
	test.NotNil(t,conf.ConfigSet)
	test.Equal(t,4,	conf.ConfigSet.Len())
}