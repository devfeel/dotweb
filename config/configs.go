package config

import (
	"encoding/xml"
	"errors"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/framework/file"
	"io/ioutil"
)

type (

	// Config dotweb app config define
	Config struct {
		XMLName        xml.Name          `xml:"config" json:"-" yaml:"-"`
		App            *AppNode          `xml:"app"`
		ConfigSetNodes []*ConfigSetNode  `xml:"configset>set"`
		Offline        *OfflineNode      `xml:"offline"`
		Server         *ServerNode       `xml:"server"`
		Session        *SessionNode      `xml:"session"`
		Routers        []*RouterNode     `xml:"routers>router"`
		Groups         []*GroupNode      `xml:"groups>group"`
		Middlewares    []*MiddlewareNode `xml:"middlewares>middleware"`
		ConfigSet      core.ReadonlyMap  `json:"-" yaml:"-"`
	}

	// OfflineNode dotweb app offline config
	OfflineNode struct {
		Offline     bool   `xml:"offline,attr"`     //是否维护，默认false
		OfflineText string `xml:"offlinetext,attr"` //当设置为维护，默认显示内容，如果设置url，优先url
		OfflineUrl  string `xml:"offlineurl,attr"`  //当设置为维护，默认维护页地址，如果设置url，优先url
	}

	// AppNode dotweb app global config
	AppNode struct {
		LogPath      string `xml:"logpath,attr"`      //文件方式日志目录，如果为空，默认当前目录
		EnabledLog   bool   `xml:"enabledlog,attr"`   //是否启用日志记录
		RunMode      string `xml:"runmode,attr"`      //运行模式，目前支持development、production
		PProfPort    int    `xml:"pprofport,attr"`    //pprof-server 端口，不能与主Server端口相同
		EnabledPProf bool   `xml:"enabledpprof,attr"` //是否启用pprof server，默认不启用
	}

	// ServerNode dotweb app's httpserver config
	ServerNode struct {
		EnabledListDir           bool   `xml:"enabledlistdir,attr"`           //设置是否启用目录浏览，仅对Router.ServerFile有效，若设置该项，则可以浏览目录文件，默认不开启
		EnabledRequestID         bool   `xml:"enabledrequestid,attr"`         //设置是否启用唯一请求ID，默认不开启，开启后使用32位UUID
		EnabledGzip              bool   `xml:"enabledgzip,attr"`              //是否启用gzip
		EnabledAutoHEAD          bool   `xml:"enabledautohead,attr"`          //设置是否自动启用Head路由，若设置该项，则会为除Websocket\HEAD外所有路由方式默认添加HEAD路由，默认不开启
		EnabledAutoCORS          bool   `xml:"enabledautocors,attr"`          //设置是否自动跨域支持，若设置，默认“GET, POST, PUT, DELETE, OPTIONS”全部请求均支持跨域
		EnabledIgnoreFavicon     bool   `xml:"enabledignorefavicon,attr"`     //设置是否忽略favicon.ico请求，若设置，网站将把所有favicon.ico请求直接空返回
		EnabledBindUseJsonTag    bool   `xml:"enabledbindusejsontag,attr"`    //设置bind是否启用json标签,默认不启用,若设置，bind自动识别json tag,忽略form tag
		Port                     int    `xml:"port,attr"`                     //端口
		EnabledTLS               bool   `xml:"enabledtls,attr"`               //是否启用TLS模式
		TLSCertFile              string `xml:"tlscertfile,attr"`              //TLS模式下Certificate证书文件地址
		TLSKeyFile               string `xml:"tlskeyfile,attr"`               //TLS模式下秘钥文件地址
		IndexPage                string `xml:"indexpage,attr"`                //默认index页面
		EnabledDetailRequestData bool   `xml:"enableddetailrequestdata,attr"` //设置状态数据是否启用详细页面统计，默认不启用，请特别对待，如果站点url过多，会导致数据量过大
	}

	// SessionNode dotweb app's session config
	SessionNode struct {
		EnabledSession bool   `xml:"enabled,attr"`  //启用Session
		SessionMode    string `xml:"mode,attr"`     //session mode，now support runtime、redis
		CookieName	   string `xml:"cookiename,attr"` //custom cookie name which sessionid store, default is dotweb_sessionId
		Timeout        int64  `xml:"timeout,attr"`  //session time-out period, with minute
		ServerIP       string `xml:"serverip,attr"` //remote session server url
		StoreKeyPre 	string `xml:"storekeypre,attr"` //remote session StoreKeyPre
	}

	// RouterNode dotweb app's router config
	RouterNode struct {
		Method      string            `xml:"method,attr"`
		Path        string            `xml:"path,attr"`
		HandlerName string            `xml:"handler,attr"`
		Middlewares []*MiddlewareNode `xml:"middleware"`
		IsUse       bool              `xml:"isuse,attr"` //是否启用，默认false
	}

	// GroupNode dotweb app's group router config
	GroupNode struct {
		Path        string            `xml:"path,attr"`
		Routers     []*RouterNode     `xml:"router"`
		Middlewares []*MiddlewareNode `xml:"middleware"`
		IsUse       bool              `xml:"isuse,attr"` //是否启用，默认false
	}

	// MiddlewareNode dotweb app's middleware config
	MiddlewareNode struct {
		Name  string `xml:"name,attr"`
		IsUse bool   `xml:"isuse,attr"` //是否启用，默认false
	}
)

const (
	// ConfigType_XML xml config file
	ConfigType_XML  = "xml"
	// ConfigType_JSON json config file
	ConfigType_JSON = "json"
	// ConfigType_Yaml yaml config file
	ConfigType_Yaml = "yaml"
)

// NewConfig create new config
func NewConfig() *Config {
	return &Config{
		App:       NewAppNode(),
		Offline:   NewOfflineNode(),
		Server:    NewServerNode(),
		Session:   NewSessionNode(),
		ConfigSet: core.NewReadonlyMap(),
	}
}

// IncludeConfigSet include ConfigSet file to Dotweb.Config.ConfigSet, can use ctx.ConfigSet to use your data
// same key will cover oldest value
// support xml\json\yaml
func (conf *Config) IncludeConfigSet(configFile string, confType string) error {
	var parseItem core.ConcurrenceMap
	var err error
	if confType == ConfigType_XML {
		parseItem, err = ParseConfigSetXML(configFile)
	}
	if confType == ConfigType_JSON {
		parseItem, err = ParseConfigSetJSON(configFile)
	}
	if confType == ConfigType_Yaml {
		parseItem, err = ParseConfigSetYaml(configFile)
	}
	if err != nil {
		return err
	}
	items := conf.ConfigSet.(*core.ItemMap)
	if items == nil {
		return errors.New("init config items error")
	}
	for k, v := range parseItem.GetCurrentMap() {
		items.Set(k, v)
	}
	return nil
}

func NewAppNode() *AppNode {
	config := &AppNode{}
	return config
}

func NewOfflineNode() *OfflineNode {
	config := &OfflineNode{}
	return config
}

func NewServerNode() *ServerNode {
	config := &ServerNode{}
	return config
}

func NewSessionNode() *SessionNode {
	config := &SessionNode{}
	return config
}

//init config file
//If an exception occurs, will be panic it
func MustInitConfig(configFile string, confType ...interface{}) *Config {
	conf, err := InitConfig(configFile, confType...)
	if err != nil {
		panic(err)
	}
	return conf
}

//初始化配置文件
//如果发生异常，返回异常
func InitConfig(configFile string, confType ...interface{}) (config *Config, err error) {

	//检查配置文件有效性
	//1、按绝对路径检查
	//2、尝试在当前进程根目录下寻找
	//3、尝试在当前进程根目录/config/ 下寻找
	//fixed for issue #15 读取配置文件路径
	realFile := configFile
	if !file.Exist(realFile) {
		realFile = file.GetCurrentDirectory() + "/" + configFile
		if !file.Exist(realFile) {
			realFile = file.GetCurrentDirectory() + "/config/" + configFile
			if !file.Exist(realFile) {
				return nil, errors.New("no exists config file => " + configFile)
			}
		}
	}

	cType := ConfigType_XML
	if len(confType) > 0 && confType[0] == ConfigType_JSON {
		cType = ConfigType_JSON
	}
	if len(confType) > 0 && confType[0] == ConfigType_Yaml {
		cType = ConfigType_Yaml
	}

	if cType == ConfigType_XML {
		config, err = initConfig(realFile, cType, UnmarshalXML)
	} else if cType == ConfigType_Yaml {
		config, err = initConfig(realFile, cType, UnmarshalYaml)
	} else {
		config, err = initConfig(realFile, cType, UnmarshalJSON)
	}

	if err != nil {
		return config, err
	}

	if config.App == nil {
		config.App = NewAppNode()
	}

	if config.Server == nil {
		config.Server = NewServerNode()
	}

	if config.Session == nil {
		config.Session = NewSessionNode()
	}

	if config.Offline == nil {
		config.Offline = NewOfflineNode()
	}

	tmpConfigSetMap := core.NewConcurrenceMap()
	for _, v := range config.ConfigSetNodes {
		tmpConfigSetMap.Set(v.Key, v.Value)
	}
	config.ConfigSet = tmpConfigSetMap

	//deal config default value
	dealConfigDefaultSet(config)

	return config, nil
}

func dealConfigDefaultSet(c *Config) {

}

func initConfig(configFile string, ctType string, parser func([]byte, interface{}) error) (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("DotWeb:Config:initConfig 当前cType:" + ctType + " 配置文件[" + configFile + "]无法解析 - " + err.Error())
	}

	var config *Config
	err = parser(content, &config)
	if err != nil {
		return nil, errors.New("DotWeb:Config:initConfig 当前cType:" + ctType + " 配置文件[" + configFile + "]解析失败 - " + err.Error())
	}
	return config, nil
}
