package servers

import "net/http"

type Server interface {
	//ServeHTTP make sure request can be handled correctly
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	//SetOffline set server offline config
	SetOffline(offline bool, offlineText string, offlineUrl string)
	//IsOffline check server is set offline state
	IsOffline() bool
}
