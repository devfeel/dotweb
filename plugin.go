package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/config"
	"os"
	"path/filepath"
	"time"
)

// Plugin a interface for app's global plugin
type Plugin interface {
	Name() string
	Run() error
	IsValidate() bool
}

// NewDefaultNotifyPlugin return new NotifyPlugin with default config
func NewDefaultNotifyPlugin(app *DotWeb) *NotifyPlugin {
	p := new(NotifyPlugin)
	p.app = app
	p.LoopTime = notifyPlugin_LoopTime
	p.Root = app.Config.ConfigFilePath
	p.suffix = make(map[string]bool)
	p.ModTimes = make(map[string]time.Time)
	return p
}

// NewNotifyPlugin return new NotifyPlugin with fileRoot & loopTime & suffix
// if suffix is nil or suffix[0] == "*", will visit all files in fileRoot
/*func NewNotifyPlugin(app *DotWeb, fileRoot string, loopTime int, suffix []string) *NotifyPlugin{
	p := new(NotifyPlugin)
	p.app = app
	p.LoopTime = loopTime
	p.Root = fileRoot
	Suffix := make(map[string]bool)
	if len(suffix) > 0 && suffix[0] != "*" {
		for _, v := range suffix {
			Suffix[v] = true
		}
	}
	p.suffix = Suffix
	p.ModTimes = make(map[string]time.Time)
	return p
}*/

const notifyPlugin_LoopTime = 500 //ms

type NotifyPlugin struct {
	app      *DotWeb
	Root     string
	suffix   map[string]bool
	LoopTime int
	ModTimes map[string]time.Time
}

func (p *NotifyPlugin) Name() string {
	return "NotifyPlugin"
}

func (p *NotifyPlugin) IsValidate() bool {
	return true
}

func (p *NotifyPlugin) Run() error {
	return p.start()
}

func (p *NotifyPlugin) visit(path string, fileinfo os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("访问文件失败%s", err)
	}
	ext := filepath.Ext(path)
	if !fileinfo.IsDir() && (p.suffix[ext] || len(p.suffix) == 0) {
		modTime := fileinfo.ModTime()
		if oldModTime, ok := p.ModTimes[path]; !ok {
			p.ModTimes[path] = modTime
		} else {
			if oldModTime.Before(modTime) {
				p.app.Logger().Info("NotifyPlugin Reload "+path, LogTarget_HttpServer)
				appConfig, err := config.InitConfig(p.app.Config.ConfigFilePath, p.app.Config.ConfigType)
				if err != nil {
					p.app.Logger().Error("NotifyPlugin Reload "+path+" error => "+fmt.Sprint(err), LogTarget_HttpServer)
				}
				p.app.ReSetConfig(appConfig)
				p.ModTimes[path] = modTime
			}
		}
	}
	return nil
}

func (p *NotifyPlugin) start() error {
	for {
		filepath.Walk(p.Root, p.visit)
		time.Sleep(time.Duration(p.LoopTime) * time.Millisecond)
	}
}
