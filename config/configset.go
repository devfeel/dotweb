package config

import (
	"encoding/xml"
	"errors"
	"github.com/devfeel/dotweb/core"
	"io/ioutil"
)

type (
	// ConfigSet 单元配置组，包含一系列单元配置节点
	ConfigSet struct {
		XMLName        xml.Name         `xml:"config" json:"-" yaml:"-"`
		Name           string           `xml:"name,attr"`
		ConfigSetNodes []*ConfigSetNode `xml:"set"`
	}

	// ConfigSetNode update for issue #16 配置文件
	ConfigSetNode struct {
		Key   string `xml:"key,attr"`
		Value string `xml:"value,attr"`
	}
)

// ParseConfigSetXML include ConfigSet xml file
func ParseConfigSetXML(configFile string) (core.ConcurrenceMap, error) {
	return parseConfigSetFile(configFile, ConfigType_XML)
}

// ParseConfigSetJSON include ConfigSet json file
func ParseConfigSetJSON(configFile string) (core.ConcurrenceMap, error) {
	return parseConfigSetFile(configFile, ConfigType_JSON)
}

// ParseConfigSetYaml include ConfigSet yaml file
func ParseConfigSetYaml(configFile string) (core.ConcurrenceMap, error) {
	return parseConfigSetFile(configFile, ConfigType_Yaml)
}

func parseConfigSetFile(configFile string, confType string) (core.ConcurrenceMap, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("DotWeb:Config:parseConfigSetFile 配置文件[" + configFile + ", " + confType + "]无法解析 - " + err.Error())
	}
	set := new(ConfigSet)
	if confType == ConfigType_XML {
		err = UnmarshalXML(content, set)
	}
	if confType == ConfigType_JSON {
		err = UnmarshalJSON(content, set)
	}
	if confType == ConfigType_Yaml {
		err = UnmarshalYaml(content, set)
	}
	if err != nil {
		return nil, errors.New("DotWeb:Config:parseConfigSetFile 配置文件[" + configFile + ", " + confType + "]无法解析 - " + err.Error())
	}
	item := core.NewConcurrenceMap()
	for _, s := range set.ConfigSetNodes {
		item.Set(s.Key, s.Value)
	}
	return item, nil
}
