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

// ServeHTTP makes the httprouter implement the http.Handler interface.
func (server *OfflineServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// maintenance mode
	if server.offline {
		// prefer url
		if server.offlineUrl != "" {
			http.Redirect(w, req, server.offlineUrl, http.StatusMovedPermanently)
		} else {
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
