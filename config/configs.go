package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type (
	AppConfig struct {
		XMLName xml.Name       `xml:"config"`
		Server  ServerConfig   `xml:"server"`
		Session SessionConfig  `xml:"session"`
		Routers []RouterConfig `xml:"routers>router"`
	}

	ServerConfig struct {
		LogPath      string `xml:"logpath,attr"`      //文件方式日志目录，如果为空，默认当前目录
		Port         int    `xml:"port,attr"`         //端口
		Offline      bool   `xml:"offline,attr"`      //是否维护，默认false
		OfflineUrl   string `xml:"offlineurl,attr"`   //当设置为维护，默认维护页地址
		EnabledDebug bool   `xml:"enableddebug,attr"` //启用Debug模式
		EnabledGzip  bool   `xml:"enabledgzip,attr"`  //启用gzip
	}

	SessionConfig struct {
		EnabledSession bool   `xml:"enabled,attr"`  //启用Session
		SessionMode    string `xml:"mode,attr"`     //session模式，目前支持runtime、redis
		Timeout        int64  `xml:"timeout,attr"`  //session超时时间，分为单位
		ServerIP       string `xml:"serverip,attr"` //远程session serverip
		UserName       string `xml:"username,attr"` //远程session username
		Password       string `xml:"password,attr"` //远程session password
	}

	RouterConfig struct {
		Method      string `xml:"method,attr"`
		Path        string `xml:"path,attr"`
		HandlerName string `xml:"handler,attr"`
		IsUse       bool   `xml:"isuse,attr"` //是否启用，默认false
	}
)

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("DotWeb:Config:InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config AppConfig
	err = xml.Unmarshal(content, &config)
	if err != nil {
		panic("DotWeb:Config:InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	return &config
}
