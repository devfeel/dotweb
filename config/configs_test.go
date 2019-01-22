package config

import (
	"testing"

	"github.com/devfeel/dotweb/test"
)

func TestInitConfig(t *testing.T) {
	conf, err := InitConfig("testdata/dotweb.json", "json")

	test.Nil(t, err)
	test.NotNil(t, conf)
	test.NotNil(t, conf.App)
	test.NotNil(t, conf.App.LogPath)
	test.NotNil(t, conf.ConfigSet)
}

func TestInitConfigWithXml(t *testing.T) {
	conf, err := InitConfig("testdata/dotweb.conf", "xml")

	test.Nil(t, err)
	test.NotNil(t, conf)
	test.NotNil(t, conf.App)
	test.NotNil(t, conf.App.LogPath)
	test.NotNil(t, conf.ConfigSet)
	//	test.Equal(t, 4, conf.ConfigSet.Len())
}
