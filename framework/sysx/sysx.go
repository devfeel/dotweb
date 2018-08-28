package sysx

import "os"

// GetHostName get host name
func GetHostName() string{
	host, _ := os.Hostname()
	return host
}
