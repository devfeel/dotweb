package core

import (
	"sync/atomic"
	"time"
)

var GlobalState *ServerStateInfo

func init() {
	GlobalState = &ServerStateInfo{
		ServerStartTime:   time.Now(),
		TotalRequestCount: 0,
		TotalErrorCount:   0,
	}
}

//服务器状态信息
type ServerStateInfo struct {
	//服务启动时间
	ServerStartTime time.Time
	//该运行期间总访问次数
	TotalRequestCount uint64
	//该运行期间错误次数
	TotalErrorCount uint64
}

//增加请求数
func (state *ServerStateInfo) AddRequestCount(num uint64) uint64 {
	atomic.AddUint64(&state.TotalRequestCount, num)
	return state.TotalRequestCount
}

//增加错误数
func (state *ServerStateInfo) AddErrorCount(num uint64) uint64 {
	atomic.AddUint64(&state.TotalErrorCount, num)
	return state.TotalErrorCount
}
