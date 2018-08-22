package dotweb

const(
	requestHeader_UseMock_Key = "dotweb_req_mock"
	requestHeader_UseMock_Flag = "true"
)


type MockHandle func(ctx Context)

type Mock interface{
	// Register reg mock handler on route
	Register(route string, handler MockHandle)
	// RegisterString reg mock return string on route
	RegisterString(route string, resData string)
	// CheckNeedMock check is need do mock logic
	CheckNeedMock(Context) bool
	// Do do mock logic
	Do(Context)
}

type StandardMock struct{
	routeMap map[string]MockHandle
}

func NewStandardMock() *StandardMock{
	return &StandardMock{routeMap:make(map[string]MockHandle)}
}

func (m *StandardMock) CheckNeedMock(ctx Context) bool{
	if ctx.Request().QueryHeader(requestHeader_UseMock_Key) == requestHeader_UseMock_Flag{
		return true
	}
	return false
}

func (m *StandardMock) Do(ctx Context){
	handler, exists:=m.routeMap[ctx.RouterNode().Node().fullPath]
	if exists{
		handler(ctx)
	}
}

func (m *StandardMock) Register(route string, handler MockHandle){
	m.routeMap[route] = handler
}

func (m *StandardMock) RegisterString(route string, resData string){
	m.routeMap[route] = func(ctx Context) {
		ctx.WriteString(resData)
		ctx.Response().SetHeader(requestHeader_UseMock_Key, requestHeader_UseMock_Flag)
		ctx.End()
	}
}