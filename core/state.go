package core

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/framework/sysx"
)

var GlobalState *ServerStateInfo

const (
	minuteTimeLayout        = "200601021504"
	dateTimeLayout          = "2006-01-02 15:04:05"
	defaultReserveMinutes   = 60
	defaultCheckTimeMinutes = 10
)

func init() {
	GlobalState = &ServerStateInfo{
		ServerStartTime:      time.Now(),
		TotalRequestCount:    0,
		TotalErrorCount:      0,
		CurrentRequestCount:  0,
		IntervalRequestData:  NewItemMap(),
		DetailRequestURLData: NewItemMap(),
		IntervalErrorData:    NewItemMap(),
		DetailErrorPageData:  NewItemMap(),
		DetailErrorData:      NewItemMap(),
		DetailHTTPCodeData:   NewItemMap(),
		dataChan_Request:     make(chan *RequestInfo, 2000),
		dataChan_Error:       make(chan *ErrorInfo, 1000),
		infoPool: &pool{
			requestInfo: sync.Pool{
				New: func() interface{} {
					return &RequestInfo{}
				},
			},
			errorInfo: sync.Pool{
				New: func() interface{} {
					return &ErrorInfo{}
				},
			},
		},
	}
	go GlobalState.handleInfo()
	go time.AfterFunc(time.Duration(defaultCheckTimeMinutes)*time.Minute, GlobalState.checkAndRemoveIntervalData)
}

type pool struct {
	requestInfo  sync.Pool
	errorInfo    sync.Pool
	httpCodeInfo sync.Pool
}

// http request count info
type RequestInfo struct {
	URL  string
	Code int
	Num  uint64
}

// error count info
type ErrorInfo struct {
	URL    string
	ErrMsg string
	Num    uint64
}

// Server state
type ServerStateInfo struct {
	ServerStartTime time.Time
	// enable detailed request statistics, fixes #63 request statistics, high memory usage when URL number is high
	EnabledDetailRequestData bool
	TotalRequestCount        uint64
	// active request count
	CurrentRequestCount uint64
	// request statistics per minute
	IntervalRequestData *ItemMap
	// detailed request statistics, the key is url without parameters
	DetailRequestURLData *ItemMap
	TotalErrorCount      uint64
	// request error statistics per minute
	IntervalErrorData *ItemMap
	// detailed request error statistics, the key is url without parameters
	DetailErrorPageData *ItemMap
	// detailed error statistics, the key is url without parameters
	DetailErrorData *ItemMap
	// detailed reponse statistics of http code, the key is HttpCode, e.g. 200, 500 etc.
	DetailHTTPCodeData *ItemMap

	dataChan_Request chan *RequestInfo
	dataChan_Error   chan *ErrorInfo
	infoPool         *pool
}

// ShowHtmlData show server state data html-string format
func (state *ServerStateInfo) ShowHtmlData(version, globalUniqueId string) string {
	data := "<html><body><div>"
	data += "GlobalUniqueId : " + globalUniqueId
	data += "<br>"
	data += "HostInfo : " + sysx.GetHostName()
	data += "<br>"
	data += "CurrentTime : " + time.Now().Format("2006-01-02 15:04:05")
	data += "<br>"
	data += "ServerVersion : " + version
	data += "<br>"
	data += "ServerStartTime : " + state.ServerStartTime.Format(dateTimeLayout)
	data += "<br>"
	data += "TotalRequestCount : " + strconv.FormatUint(state.TotalRequestCount, 10)
	data += "<br>"
	data += "CurrentRequestCount : " + strconv.FormatUint(state.CurrentRequestCount, 10)
	data += "<br>"
	data += "TotalErrorCount : " + strconv.FormatUint(state.TotalErrorCount, 10)
	data += "<br>"
	state.IntervalRequestData.RLock()
	data += "IntervalRequestData : " + jsonutil.GetJsonString(state.IntervalRequestData.GetCurrentMap())
	state.IntervalRequestData.RUnlock()
	data += "<br>"
	state.DetailRequestURLData.RLock()
	data += "DetailRequestUrlData : " + jsonutil.GetJsonString(state.DetailRequestURLData.GetCurrentMap())
	state.DetailRequestURLData.RUnlock()
	data += "<br>"
	state.IntervalErrorData.RLock()
	data += "IntervalErrorData : " + jsonutil.GetJsonString(state.IntervalErrorData.GetCurrentMap())
	state.IntervalErrorData.RUnlock()
	data += "<br>"
	state.DetailErrorPageData.RLock()
	data += "DetailErrorPageData : " + jsonutil.GetJsonString(state.DetailErrorPageData.GetCurrentMap())
	state.DetailErrorPageData.RUnlock()
	data += "<br>"
	state.DetailErrorData.RLock()
	data += "DetailErrorData : " + jsonutil.GetJsonString(state.DetailErrorData.GetCurrentMap())
	state.DetailErrorData.RUnlock()
	data += "<br>"
	state.DetailHTTPCodeData.RLock()
	data += "DetailHttpCodeData : " + jsonutil.GetJsonString(state.DetailHTTPCodeData.GetCurrentMap())
	state.DetailHTTPCodeData.RUnlock()
	data += "</div></body></html>"
	return data
}

// QueryIntervalRequestData query request count by query time
func (state *ServerStateInfo) QueryIntervalRequestData(queryKey string) uint64 {
	return state.IntervalRequestData.GetUInt64(queryKey)
}

// QueryIntervalErrorData query error count by query time
func (state *ServerStateInfo) QueryIntervalErrorData(queryKey string) uint64 {
	return state.IntervalErrorData.GetUInt64(queryKey)
}

// AddRequestCount add request count
func (state *ServerStateInfo) AddRequestCount(page string, code int, num uint64) {
	state.addRequestData(page, code, num)
}

// AddCurrentRequest increment current request count
func (state *ServerStateInfo) AddCurrentRequest(num uint64) uint64 {
	atomic.AddUint64(&state.CurrentRequestCount, num)
	return state.CurrentRequestCount
}

// SubCurrentRequest subtract current request count
func (state *ServerStateInfo) SubCurrentRequest(num uint64) uint64 {
	atomic.AddUint64(&state.CurrentRequestCount, ^uint64(num-1))
	return state.CurrentRequestCount
}

// AddErrorCount add error count
func (state *ServerStateInfo) AddErrorCount(page string, err error, num uint64) uint64 {
	atomic.AddUint64(&state.TotalErrorCount, num)
	state.addErrorData(page, err, num)
	return state.TotalErrorCount
}

func (state *ServerStateInfo) addRequestData(page string, code int, num uint64) {
	// get from pool
	info := state.infoPool.requestInfo.Get().(*RequestInfo)
	info.URL = page
	info.Code = code
	info.Num = num
	state.dataChan_Request <- info
}

func (state *ServerStateInfo) addErrorData(page string, err error, num uint64) {
	// get from pool
	info := state.infoPool.errorInfo.Get().(*ErrorInfo)
	info.URL = page
	info.ErrMsg = err.Error()
	info.Num = num
	state.dataChan_Error <- info
}

// handle logging
func (state *ServerStateInfo) handleInfo() {
	for {
		select {
		case info := <-state.dataChan_Request:
			{
				if strings.Index(info.URL, "/dotweb/") != 0 {
					atomic.AddUint64(&state.TotalRequestCount, info.Num)
				}
				// fixes #63 request statistics, high memory usage when URL number is high
				if state.EnabledDetailRequestData {
					// ignore 404 request
					if info.Code != http.StatusNotFound {
						// set detail url data
						key := strings.ToLower(info.URL)
						val := state.DetailRequestURLData.GetUInt64(key)
						state.DetailRequestURLData.Set(key, val+info.Num)
					}
				}
				// set interval data
				key := time.Now().Format(minuteTimeLayout)
				val := state.IntervalRequestData.GetUInt64(key)
				state.IntervalRequestData.Set(key, val+info.Num)

				// set code data
				key = strconv.Itoa(info.Code)
				val = state.DetailHTTPCodeData.GetUInt64(key)
				state.DetailHTTPCodeData.Set(key, val+info.Num)

				// put info obj
				state.infoPool.requestInfo.Put(info)
			}
		case info := <-state.dataChan_Error:
			{
				// set detail error page data
				key := strings.ToLower(info.URL)
				val := state.DetailErrorPageData.GetUInt64(key)
				state.DetailErrorPageData.Set(key, val+info.Num)

				// set detail error data
				key = info.ErrMsg
				val = state.DetailErrorData.GetUInt64(key)
				state.DetailErrorData.Set(key, val+info.Num)

				// set interval data
				key = time.Now().Format(minuteTimeLayout)
				val = state.IntervalErrorData.GetUInt64(key)
				state.IntervalErrorData.Set(key, val+info.Num)

				// put info obj
				state.infoPool.errorInfo.Put(info)
			}
		}
	}
}

// check and remove need to remove interval data with request and error
func (state *ServerStateInfo) checkAndRemoveIntervalData() {
	var needRemoveKey []string
	now, _ := time.Parse(minuteTimeLayout, time.Now().Format(minuteTimeLayout))

	// check IntervalRequestData
	state.IntervalRequestData.RLock()
	if state.IntervalRequestData.Len() > defaultReserveMinutes {
		for k := range state.IntervalRequestData.GetCurrentMap() {
			if t, err := time.Parse(minuteTimeLayout, k); err != nil {
				needRemoveKey = append(needRemoveKey, k)
			} else {
				if now.Sub(t) > (defaultReserveMinutes * time.Minute) {
					needRemoveKey = append(needRemoveKey, k)
				}
			}
		}
	}
	state.IntervalRequestData.RUnlock()
	// remove keys
	for _, v := range needRemoveKey {
		state.IntervalRequestData.Remove(v)
	}

	// check IntervalErrorData
	needRemoveKey = []string{}
	state.IntervalErrorData.RLock()
	if state.IntervalErrorData.Len() > defaultReserveMinutes {
		for k := range state.IntervalErrorData.GetCurrentMap() {
			if t, err := time.Parse(minuteTimeLayout, k); err != nil {
				needRemoveKey = append(needRemoveKey, k)
			} else {
				if now.Sub(t) > (defaultReserveMinutes * time.Minute) {
					needRemoveKey = append(needRemoveKey, k)
				}
			}
		}
	}
	state.IntervalErrorData.RUnlock()
	// remove keys
	for _, v := range needRemoveKey {
		state.IntervalErrorData.Remove(v)
	}
	time.AfterFunc(time.Duration(defaultCheckTimeMinutes)*time.Minute, state.checkAndRemoveIntervalData)
}
