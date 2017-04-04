package config

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type (
	Config struct {
		XMLName xml.Name        `xml:"config"`
		App     *AppConfig      `xml:"app"`
		Offline *OfflineConfig  `xml:"offline"`
		Server  *ServerConfig   `xml:"server"`
		Session *SessionConfig  `xml:"session"`
		Routers []*RouterConfig `xml:"routers>router"`
	}
	OfflineConfig struct {
		Offline     bool   `xml:"offline,attr"`     //是否维护，默认false
		OfflineText string `xml:"offlinetext,attr"` //当设置为维护，默认显示内容，如果设置url，优先url
		OfflineUrl  string `xml:"offlineurl,attr"`  //当设置为维护，默认维护页地址，如果设置url，优先url
	}
	AppConfig struct {
		LogPath    string `xml:"logpath,attr"`    //文件方式日志目录，如果为空，默认当前目录
		EnabledLog bool   `xml:"enabledlog,attr"` //是否启用日志记录
		RunMode    string `xml:"runmode,attr"`    //运行模式，目前支持development、production
	}
	ServerConfig struct {
		EnabledListDir  bool `xml:"enabledlistdir,attr"`  //设置是否启用目录浏览，仅对Router.ServerFile有效，若设置该项，则可以浏览目录文件，默认不开启
		EnabledGzip     bool `xml:"enabledgzip,attr"`     //是否启用gzip
		EnabledAutoHEAD bool `xml:"enabledautohead,attr"` //设置是否自动启用Head路由，若设置该项，则会为除Websocket\HEAD外所有路由方式默认添加HEAD路由，默认不开启
		Port            int  `xml:"port,attr"`            //端口
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

func NewAppConfig() *AppConfig {
	config := &AppConfig{}
	return config
}
func NewServerConfig() *ServerConfig {
	config := &ServerConfig{}
	return config
}

func NewSessionConfig() *SessionConfig {
	config := &SessionConfig{}
	return config
}

//初始化配置文件
func InitConfig(configFile string) *Config {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("DotWeb:Config:InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config Config
	err = xml.Unmarshal(content, &config)
	if err != nil {
		panic("DotWeb:Config:InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}

	if config.App == nil {
		config.App = &AppConfig{}
	}

	if config.Server == nil {
		config.Server = &ServerConfig{}
	}

	if config.Session == nil {
		config.Session = &SessionConfig{}
	}

	if config.Offline == nil {
		config.Offline = &OfflineConfig{}
	}
	return &config
}
