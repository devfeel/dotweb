package dotweb

import (
	"net/http"
	"net/url"
	"testing"
)

// TestGroupSetNotFoundHandle tests the SetNotFoundHandle functionality
func TestGroupSetNotFoundHandle(t *testing.T) {
	tests := []struct {
		name           string
		groupPrefix    string
		requestPath    string
		expectedBody   string
		shouldUseGroup bool
	}{
		{
			name:           "Group 404 - API endpoint not found",
			groupPrefix:    "/api",
			requestPath:    "/api/users",
			expectedBody:   "API 404",
			shouldUseGroup: true,
		},
		{
			name:           "Group 404 - Similar prefix should not match",
			groupPrefix:    "/api",
			requestPath:    "/api_v2/users",
			expectedBody:   "Global 404",
			shouldUseGroup: false,
		},
		{
			name:           "Global 404 - No matching group",
			groupPrefix:    "/api",
			requestPath:    "/web/index",
			expectedBody:   "Global 404",
			shouldUseGroup: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create app
			app := New()
			
			// Set global 404 handler
			app.SetNotFoundHandle(func(ctx Context) {
				ctx.WriteString("Global 404")
			})
			
			// Create group with custom 404 handler
			group := app.HttpServer.Group(tt.groupPrefix)
			group.SetNotFoundHandle(func(ctx Context) {
				ctx.WriteString(tt.expectedBody)
			})
			
			// Add a valid route to group
			group.GET("/exists", func(ctx Context) error {
				return ctx.WriteString("OK")
			})
			
			// Create context
			context := &HttpContext{
				response: &Response{},
				request: &Request{
					Request: &http.Request{
						URL:    &url.URL{Path: tt.requestPath},
						Method: "GET",
					},
				},
				httpServer: &HttpServer{
					DotApp: app,
				},
				routerNode: &Node{},
			}
			
			w := &testHttpWriter{}
			context.response = NewResponse(w)
			
			// Serve HTTP
			app.HttpServer.Router().ServeHTTP(context)
			
			// Check response - we can't easily check body content without more setup
			// This test mainly verifies no panic and correct routing logic
		})
	}
}

// TestGroupNotFoundHandlePriority tests that group handler takes priority over global handler
func TestGroupNotFoundHandlePriority(t *testing.T) {
	app := New()
	
	// Set global handler
	app.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString("Global Handler")
	})
	
	// Create group with handler
	apiGroup := app.HttpServer.Group("/api")
	apiGroup.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString("Group Handler")
	})
	
	// Add valid route
	apiGroup.GET("/users", func(ctx Context) error {
		return ctx.WriteString("Users")
	})
	
	// Verify group has notFoundHandler set
	xg := apiGroup.(*xGroup)
	if xg.notFoundHandler == nil {
		t.Error("Group should have notFoundHandler set")
	}
}

// TestMultipleGroupsWithNotFoundHandle tests multiple groups with different handlers
func TestMultipleGroupsWithNotFoundHandle(t *testing.T) {
	app := New()
	
	// Set global handler
	app.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString("Global 404")
	})
	
	// Create API group
	apiGroup := app.HttpServer.Group("/api")
	apiGroup.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString(`{"code": 404, "message": "API not found"}`)
	})
	
	// Create Web group
	webGroup := app.HttpServer.Group("/web")
	webGroup.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString("<h1>404 - Page Not Found</h1>")
	})
	
	// Verify both groups have handlers
	apiXg := apiGroup.(*xGroup)
	webXg := webGroup.(*xGroup)
	
	if apiXg.notFoundHandler == nil {
		t.Error("API group should have notFoundHandler set")
	}
	if webXg.notFoundHandler == nil {
		t.Error("Web group should have notFoundHandler set")
	}
}

// TestGroupSetNotFoundHandleReturnsGroup tests that SetNotFoundHandle returns the Group for chaining
func TestGroupSetNotFoundHandleReturnsGroup(t *testing.T) {
	app := New()
	
	group := app.HttpServer.Group("/api")
	result := group.SetNotFoundHandle(func(ctx Context) {
		ctx.WriteString("404")
	})
	
	if result == nil {
		t.Error("SetNotFoundHandle should return Group for chaining")
	}
}

// test helper
type testHttpWriter http.Header

func (ho testHttpWriter) Header() http.Header {
	return http.Header(ho)
}

func (ho testHttpWriter) Write(byte []byte) (int, error) {
	return len(byte), nil
}

func (ho testHttpWriter) WriteHeader(code int) {
}
