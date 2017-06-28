package config

import (
	"testing"
	"github.com/zouyx/agollo/test"
)

func TestInitConfig(t *testing.T) {
	conf,err:=InitConfig("example/config/dotweb.json.conf","json")

	test.Nil(t,err)
	test.NotNil(t,conf)
	test.NotNil(t,conf.App)
	test.NotNil(t,conf.App.LogPath)
	test.NotNil(t,conf.AppSets)
	test.Equal(t,4,	len(conf.AppSets))
}

func TestInitConfigWithXml(t *testing.T) {
	conf,err:=InitConfig("example/config/dotweb.conf","xml")

	test.Nil(t,err)
	test.NotNil(t,conf)
	test.NotNil(t,conf.App)
	test.NotNil(t,conf.App.LogPath)
	test.NotNil(t,conf.AppSets)
	test.Equal(t,4,	len(conf.AppSets))
}