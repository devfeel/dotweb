package core

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/framework/sysx"
)

const (
	minuteTimeLayout        = "200601021504"
	dateTimeLayout          = "2006-01-02 15:04:05"
	defaultReserveMinutes   = 60
	defaultCheckTimeMinutes = 10
)

var TableHtml = `<html>
<html><head>
   <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0;">
 
  <meta name="Generator" content="EditPlus">
  <meta name="Author" content="">
  <meta name="Keywords" content="">
  <meta name="Description" content="">
  <title>Dotweb</title>
    <style>
    .overtable {
      width: 100%;
      overflow: hidden;
      overflow-x: auto;
    }
    body {
      max-width: 780px;
       margin:0 auto;
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      font-size: 1rem;
      color: #444;
    }
    table {
      font-family: 'trebuchet MS', 'Lucida sans', Arial;
      *border-collapse: collapse;
      /* IE7 and lower */
      border-spacing: 0;
      width: 100%;
      border-collapse: collapse;
      overflow-x: auto
    }
    caption {
      font-family: 'Microsoft Yahei', 'trebuchet MS', 'Lucida sans', Arial;
    text-align: left;
    padding: .5rem;
    font-weight: bold;
    font-size: 110%;
    color: #666;
    }
    tr {
      border-top: 1px solid #dfe2e5
    }
    tr:nth-child(2n) {
      background-color: #f6f8fa
    }
    td,
    th {
      border: 1px solid #dfe2e5;
      padding: .6em 1em;
    }
    .bordered tr:hover {
      background: #fbf8e9;
    }
    .bordered td,
    .bordered th {
      border: 1px solid #ccc;
      padding: 10px;
      text-align: left;
    }
  </style>
  <script>
  
(function(doc, win) {
    window.MPIXEL_RATIO = (function () {
        var Mctx = document.createElement("canvas").getContext("2d"),
            Mdpr = window.devicePixelRatio || 1,
            Mbsr = Mctx.webkitBackingStorePixelRatio ||
                Mctx.mozBackingStorePixelRatio ||
                Mctx.msBackingStorePixelRatio ||
                Mctx.oBackingStorePixelRatio ||
                Mctx.backingStorePixelRatio || 1;
    
        return Mdpr/Mbsr;
    })();

    function addEventListeners(ele,type,callback){
    
        try{  // Chrome、FireFox、Opera、Safari、IE9.0及其以上版本
            ele.addEventListener(type,callback,false);
        }catch(e){
            try{  // IE8.0及其以下版本
                ele.attachEvent('on' + type,callback);
            }catch(e){  // 早期浏览器
                ele['on' + type] = callback;
            }
        }
    }

    var docEl = doc.documentElement,
        resizeEvt = 'orientationchange' in window ? 'orientationchange' : 'resize';
    window.recalc = function() {
            var clientWidth = docEl.clientWidth < 768 ? docEl.clientWidth : 768;
            if (!clientWidth) return;
            docEl.style.fontSize = 10 * (clientWidth / 320) *  window.MPIXEL_RATIO + 'px';
        };
    window.recalc();
    
    addEventListeners(win, resizeEvt, recalc);
})(document, window);

</script>
</head>
<body>
<div class="overtable">
{{tableBody}}
</div>
</body>
</html>
`

// NewServerStateInfo return ServerStateInfo which is init
func NewServerStateInfo() *ServerStateInfo {
	state := &ServerStateInfo{
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
	go state.handleInfo()
	go time.AfterFunc(time.Duration(defaultCheckTimeMinutes)*time.Minute, state.checkAndRemoveIntervalData)
	return state
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

// ShowHtmlDataRaw show server state data html-string format
func (state *ServerStateInfo) ShowHtmlDataRaw(version, globalUniqueId string) string {
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

// ShowHtmlData show server state data html-table format
func (state *ServerStateInfo) ShowHtmlTableData(version, globalUniqueId string) string {
	data := "<tr><td>" + "GlobalUniqueId" + "</td><td>" + globalUniqueId + "</td></tr>"
	data += "<tr><td>" + "HostInfo" + "</td><td>" + sysx.GetHostName() + "</td></tr>"
	data += "<tr><td>" + "CurrentTime" + "</td><td>" + time.Now().Format("2006-01-02 15:04:05") + "</td></tr>"
	data += "<tr><td>" + "ServerVersion" + "</td><td>" + version + "</td></tr>"
	data += "<tr><td>" + "ServerStartTime" + "</td><td>" + state.ServerStartTime.Format(dateTimeLayout) + "</td></tr>"
	data += "<tr><td>" + "TotalRequestCount" + "</td><td>" + strconv.FormatUint(state.TotalRequestCount, 10) + "</td></tr>"
	data += "<tr><td>" + "CurrentRequestCount" + "</td><td>" + strconv.FormatUint(state.CurrentRequestCount, 10) + "</td></tr>"
	data += "<tr><td>" + "TotalErrorCount" + "</td><td>" + strconv.FormatUint(state.TotalErrorCount, 10) + "</td></tr>"
	state.IntervalErrorData.RLock()
	data += "<tr><td>" + "IntervalErrorData" + "</td><td>" + jsonutil.GetJsonString(state.IntervalErrorData.GetCurrentMap()) + "</td></tr>"
	state.IntervalErrorData.RUnlock()
	state.DetailErrorPageData.RLock()
	data += "<tr><td>" + "DetailErrorPageData" + "</td><td>" + jsonutil.GetJsonString(state.DetailErrorPageData.GetCurrentMap()) + "</td></tr>"
	state.DetailErrorPageData.RUnlock()
	state.DetailErrorData.RLock()
	data += "<tr><td>" + "DetailErrorData" + "</td><td>" + jsonutil.GetJsonString(state.DetailErrorData.GetCurrentMap()) + "</td></tr>"
	state.DetailErrorData.RUnlock()
	state.DetailHTTPCodeData.RLock()
	data += "<tr><td>" + "DetailHttpCodeData" + "</td><td>" + jsonutil.GetJsonString(state.DetailHTTPCodeData.GetCurrentMap()) + "</td></tr>"
	state.DetailHTTPCodeData.RUnlock()
	header := `<tr>
          <th>Index</th>
          <th>Value</th>
        </tr>`
	data = CreateTableHtml("Core State", header, data)

	//show IntervalRequestData
	intervalRequestData := ""
	state.IntervalRequestData.RLock()
	for k, v := range state.IntervalRequestData.GetCurrentMap() {
		intervalRequestData += "<tr><td>" + k + "</td><td>" + fmt.Sprint(v) + "</td></tr>"
	}
	state.IntervalRequestData.RUnlock()
	header = `<tr>
          <th>Time</th>
          <th>Value</th>
        </tr>`
	data += CreateTableHtml("IntervalRequestData", header, intervalRequestData)

	//show DetailRequestURLData
	detailRequestURLData := ""
	state.DetailRequestURLData.RLock()
	for k, v := range state.DetailRequestURLData.GetCurrentMap() {
		detailRequestURLData += "<tr><td>" + k + "</td><td>" + fmt.Sprint(v) + "</td></tr>"
	}
	state.DetailRequestURLData.RUnlock()
	header = `<tr>
          <th>Url</th>
          <th>Value</th>
        </tr>`
	data += CreateTableHtml("DetailRequestURLData", header, detailRequestURLData)
	html := strings.Replace(TableHtml, "{{tableBody}}", data, -1)

	return html
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
