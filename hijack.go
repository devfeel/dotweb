package dotweb

import (
	"bufio"
	"net"
)

//hijack conn
type HijackConn struct {
	ReadWriter *bufio.ReadWriter
	Conn       net.Conn
	header     string
}

// WriteString hjiack conn write string
func (hj *HijackConn) WriteString(content string) (int, error) {
	n, err := hj.ReadWriter.WriteString(hj.header + "\r\n" + content)
	if err == nil {
		hj.ReadWriter.Flush()
	}
	return n, err
}

// WriteBlob hjiack conn write []byte
func (hj *HijackConn) WriteBlob(p []byte) (size int, err error) {
	size, err = hj.ReadWriter.Write(p)
	if err == nil {
		hj.ReadWriter.Flush()
	}
	return
}

// SetHeader hjiack conn write header
func (hj *HijackConn) SetHeader(key, value string) {
	hj.header += key + ": " + value + "\r\n"
}

// Close close hijack conn
func (hj *HijackConn) Close() error {
	return hj.Conn.Close()
}
