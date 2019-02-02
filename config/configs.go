package config

import (
	"encoding/xml"
	"errors"
	"io/ioutil"

	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/framework/file"
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
		Offline     bool   `xml:"offline,attr"`     // maintenance mode, default false
		OfflineText string `xml:"offlinetext,attr"` // text to display when Offline is true, OfflineUrl is used if set
		OfflineUrl  string `xml:"offlineurl,attr"`  // maintenance page url
	}

	// AppNode dotweb app global config
	AppNode struct {
		LogPath      string `xml:"logpath,attr"`      // path of log files, use current directory if empty
		EnabledLog   bool   `xml:"enabledlog,attr"`   // enable logging
		RunMode      string `xml:"runmode,attr"`      // run mode, currently supports [development, production]
		PProfPort    int    `xml:"pprofport,attr"`    // pprof-server port, cann't be same as server port
		EnabledPProf bool   `xml:"enabledpprof,attr"` // enable pprof server, default is false
	}

	// ServerNode dotweb app's httpserver config
	ServerNode struct {
		EnabledListDir              bool   `xml:"enabledlistdir,attr"`   // enable listing of directories, only valid for Router.ServerFile, default is false
		EnabledRequestID            bool   `xml:"enabledrequestid,attr"` // enable uniq request ID, default is false, 32-bit UUID is used if enabled
		EnabledGzip                 bool   `xml:"enabledgzip,attr"`      // enable gzip
		EnabledAutoHEAD             bool   `xml:"enabledautohead,attr"`  // ehanble HEAD routing, default is false, will add HEAD routing for all routes except for websocket and HEAD
		EnabledAutoOPTIONS          bool   // enable OPTIONS routing, default is false, will add OPTIONS routing for all routes except for websocket and OPTIONS
		EnabledIgnoreFavicon        bool   `xml:"enabledignorefavicon,attr"`  // ignore favicon.ico request, return empty reponse if set
		EnabledBindUseJsonTag       bool   `xml:"enabledbindusejsontag,attr"` // allow Bind to use JSON tag, default is false, Bind will use json tag automatically and ignore form tag
		EnabledStaticFileMiddleware bool   // The flag which enabled or disabled middleware for static-file route
		Port                        int    `xml:"port,attr"`                     // port
		EnabledTLS                  bool   `xml:"enabledtls,attr"`               // enable TLS
		TLSCertFile                 string `xml:"tlscertfile,attr"`              // certifications file for TLS
		TLSKeyFile                  string `xml:"tlskeyfile,attr"`               // keys file for TLS
		IndexPage                   string `xml:"indexpage,attr"`                // default index page
		EnabledDetailRequestData    bool   `xml:"enableddetailrequestdata,attr"` // enable detailed statics for requests, default is false. Please use with care, it will have performance issues if the site have lots of URLs
		VirtualPath                 string // virtual path when deploy on no root path
	}

	// SessionNode dotweb app's session config
	SessionNode struct {
		EnabledSession  bool   `xml:"enabled,attr"`         // enable session
		SessionMode     string `xml:"mode,attr"`            // session mode，now support runtime、redis
		CookieName      string `xml:"cookiename,attr"`      // custom cookie name which sessionid store, default is dotweb_sessionId
		Timeout         int64  `xml:"timeout,attr"`         // session time-out period, with second
		ServerIP        string `xml:"serverip,attr"`        // remote session server url
		BackupServerUrl string `xml:"backupserverurl,attr"` // backup remote session server url
		StoreKeyPre     string `xml:"storekeypre,attr"`     // remote session StoreKeyPre
	}

	// RouterNode dotweb app's router config
	RouterNode struct {
		Method      string            `xml:"method,attr"`
		Path        string            `xml:"path,attr"`
		HandlerName string            `xml:"handler,attr"`
		Middlewares []*MiddlewareNode `xml:"middleware"`
		IsUse       bool              `xml:"isuse,attr"` // enable router, default is false
	}

	// GroupNode dotweb app's group router config
	GroupNode struct {
		Path        string            `xml:"path,attr"`
		Routers     []*RouterNode     `xml:"router"`
		Middlewares []*MiddlewareNode `xml:"middleware"`
		IsUse       bool              `xml:"isuse,attr"` // enable group, default is false
	}

	// MiddlewareNode dotweb app's middleware config
	MiddlewareNode struct {
		Name  string `xml:"name,attr"`
		IsUse bool   `xml:"isuse,attr"` // enable middleware, default is false
	}
)

const (
	// ConfigType_XML xml config file
	ConfigType_XML = "xml"
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

// init config file
// If an exception occurs, will be panic it
func MustInitConfig(configFile string, confType ...interface{}) *Config {
	conf, err := InitConfig(configFile, confType...)
	if err != nil {
		panic(err)
	}
	return conf
}

// InitConfig initialize the config with configFile
func InitConfig(configFile string, confType ...interface{}) (config *Config, err error) {

	// Validity check
	// 1. Try read as absolute path
	// 2. Try the current working directory
	// 3. Try $PWD/config
	// fixed for issue #15 config file path
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

	// deal config default value
	dealConfigDefaultSet(config)

	return config, nil
}

func dealConfigDefaultSet(c *Config) {

}

func initConfig(configFile string, ctType string, parser func([]byte, interface{}) error) (*Config, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("DotWeb:Config:initConfig current cType:" + ctType + " config file [" + configFile + "] cannot be parsed - " + err.Error())
	}

	var config *Config
	err = parser(content, &config)
	if err != nil {
		return nil, errors.New("DotWeb:Config:initConfig current cType:" + ctType + " config file [" + configFile + "] cannot be parsed - " + err.Error())
	}
	return config, nil
}
