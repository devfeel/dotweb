package servers

import "net/http"

const (
	DefaultOfflineText = "sorry, server is offline!"
	NotSetOfflineText  = "why you come here?"
)

type OfflineServer struct {
	offline     bool
	offlineText string
	offlineUrl  string
}

func NewOfflineServer() Server {
	return &OfflineServer{}
}

func (server *OfflineServer) IsOffline() bool {
	return server.offline
}

func (server *OfflineServer) SetOffline(offline bool, offlineText string, offlineUrl string) {
	server.offline = offline
	server.offlineUrl = offlineUrl
	server.offlineText = offlineText
}

//ServeHTTP makes the httprouter implement the http.Handler interface.
func (server *OfflineServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//处理维护
	if server.offline {
		//url优先
		if server.offlineUrl != "" {
			http.Redirect(w, req, server.offlineUrl, http.StatusMovedPermanently)
		} else {
			//输出内容
			if server.offlineText == "" {
				server.offlineText = DefaultOfflineText
			}
			w.Write([]byte(server.offlineText))
		}
		return
	} else {
		w.Write([]byte(NotSetOfflineText))
	}
}
