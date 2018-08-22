package dotweb

const(
	requestHeaderUseMockKey = "dotweb_req_mock"
	requestHeaderUseMockFlag = "true"
)

// MockHandle the handle define on mock module
type MockHandle func(ctx Context)

// Mock the define Mock module
type Mock interface{
	// Register register MockHandle on route
	Register(route string, handler MockHandle)
	// RegisterString register return mock string on route
	RegisterString(route string, resData interface{})
	// RegisterJSON register return mock json on route
	RegisterJSON(route string, resData interface{})
	// CheckNeedMock check is need do mock logic
	CheckNeedMock(Context) bool
	// Do do mock logic
	Do(Context)
}

// StandardMock standard mock implement for Mock interface
type StandardMock struct{
	routeMap map[string]MockHandle
}

// NewStandardMock create new StandardMock
func NewStandardMock() *StandardMock{
	return &StandardMock{routeMap:make(map[string]MockHandle)}
}

// CheckNeedMock check is need do mock logic
func (m *StandardMock) CheckNeedMock(ctx Context) bool{
	if ctx.Request().QueryHeader(requestHeaderUseMockKey) == requestHeaderUseMockFlag{
		return true
	}
	return false
}

// Do do mock logic
func (m *StandardMock) Do(ctx Context){
	handler, exists:=m.routeMap[ctx.RouterNode().Node().fullPath]
	if exists{
		handler(ctx)
	}
}

// Register register MockHandle on route
func (m *StandardMock) Register(route string, handler MockHandle){
	m.routeMap[route] = handler
}

// RegisterString register return mock string on route
func (m *StandardMock) RegisterString(route string, resData interface{}){
	m.routeMap[route] = func(ctx Context) {
		ctx.Response().SetHeader(requestHeaderUseMockKey, requestHeaderUseMockFlag)
		ctx.WriteString(resData)
		ctx.End()
	}
}

// RegisterJSON register return mock json on route
func (m *StandardMock) RegisterJSON(route string, resData interface{}){
	m.routeMap[route] = func(ctx Context) {
		ctx.Response().SetHeader(requestHeaderUseMockKey, requestHeaderUseMockFlag)
		ctx.WriteJson(resData)
		ctx.End()
	}
}