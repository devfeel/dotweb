# dotweb Issues 修复分析

## Issue #245: 尾斜杠路由问题

### 问题描述
请求 `/hello/xxx/` 返回 301 重定向到 `/hello/xxx`，但期望返回 404（与 net/http 标准库一致）

### 测试代码
```go
package main

import (
    "github.com/devfeel/dotweb"
)

func main() {
    dotapp := dotweb.New()

    dotapp.HttpServer.GET("/hello/:name", func(ctx dotweb.Context) error{
        return ctx.WriteString("hello " + ctx.GetRouterName("name"))
    })

    dotapp.StartServer(80)
}
```

### 测试结果
```bash
$ curl http://localhost/hello/xx/
<a href="/hello/xx">Moved Permanently</a>

期望: 404 page not found # 与 net/http 标准库一致
```

### 根本原因分析

1. **httprouter 行为**：dotweb 基于 httprouter，httprouter 会自动清理 URL 尾部的斜杠
2. **与标准库不一致**：Go 的 net/http 标准库不会自动重定向，而是直接返回 404

### 修复方案

#### 方案 A：添加配置选项（推荐）
在 `HttpServer` 结构体中添加 `AutoRedirectTrailingSlash` 配置项：

```go
type ServerConfig struct {
    // 现有字段...
    AutoRedirectTrailingSlash bool  // 新增：是否自动重定向尾斜杠（默认 false）
}
```

**优点**：
- 向后兼容（默认 false 可以保持旧行为）
- 用户可以选择是否启用自动重定向
- 与标准库行为一致（默认 false）

**缺点**：
- 需要修改配置结构
- 需要在 Server 初始化时处理此配置

#### 方案 B：修改路由匹配逻辑
在 `router.go` 的 `MatchPath` 函数中添加特殊处理：

```go
func (r *router) MatchPath(ctx Context, routePath string) bool {
    // 现有逻辑...

    // 添加：如果请求路径以 / 结尾且匹配到无尾斜杠的路由
    // 检查是否返回 404 而不是重定向
    if strings.HasSuffix(ctx.Request().URL.Path, "/") {
        // 清理尾斜杠后重新匹配
        trimmedPath := strings.TrimSuffix(ctx.Request().URL.Path, "/")
        if r.MatchPath(ctx, trimmedPath) {
            return false // 返回 false 让路由继续匹配，最终返回 404
        }
    }

    // 原有逻辑...
}
```

**优点**：
- 不需要修改配置
- 直接修复问题
- 与标准库行为一致

**缺点**：
- 可能影响现有路由逻辑
- 需要全面测试

#### 方案 C：在中间件中处理
创建一个新的中间件来阻止自动重定向：

```go
// DisableTrailingSlashRedirect 禁用尾斜杠重定向中间件
func DisableTrailingSlashRedirect() Middleware {
    return func(ctx Context) error {
        // 设置一个标记，在路由匹配时检查
        // 如果请求路径以 / 结尾，直接清理后继续
        if strings.HasSuffix(ctx.Request().URL.Path, "/") {
            // 修改 Request URI，移除尾斜杠
            // 注意：这可能需要直接操作 http.Request
        }
        return ctx.Next()
    }
}
```

**优点**：
- 作为中间件，可以灵活注册
- 不需要修改核心路由逻辑

**缺点**：
- 实现复杂度较高
- 可能影响性能

### 推荐实施步骤

1. **方案 A（配置选项）**
   - 修改 `server.go`，添加 `AutoRedirectTrailingSlash` 配置
   - 修改 `router.go`，读取配置并处理
   - 添加测试用例

2. **测试验证**
   - 测试带尾斜杠的 URL 返回 404
   - 测试不带尾斜杠的 URL 正常工作
   - 测试动态路由参数不受影响

3. **文档更新**
   - 更新 README.md 说明新配置
   - 添加示例说明如何禁用自动重定向

---

## Issue #250: CORS 中间件问题

### 问题描述
在路由组上注册 CORS 中间件时：
```go
iweb.HttpServer.Group("/api").Use(cors.New(""))
```

**简单跨域请求**：一切正常

**复杂跨域请求** + `SetEnabledAutoOPTIONS(true)`：会提示已被 CORS 策略阻止错误。

如果手动注册 OPTIONS 请求则一切正常：
```go
apiGroup.OPTIONS("/app", dotweb.DefaultAutoOPTIONSHandler)
```

### 根本原因分析

1. **OPTIONS 自动注册时机**：`SetEnabledAutoOPTIONS(true)` 的执行时机晚于中间件注册
2. **中间件链执行顺序**：CORS 中间件可能在 OPTIONS 请求处理前执行，导致预检失败
3. **OPTIONS 请求未通过中间件链**：自动注册的 OPTIONS 处理器可能绕过了 CORS 中间件

### 修复方案

#### 方案 A：在 CORS 中间件中处理 OPTIONS（推荐）

修改 CORS 中间件，确保 OPTIONS 请求也经过 CORS 处理：

```go
// cors_middleware.go
func New(allowOrigin string) Middleware {
    return func(ctx Context) error {
        // OPTIONS 预检请求
        if ctx.Request().Method == "OPTIONS" {
            // 设置 CORS 响应头
            headers := ctx.Response().Header()
            headers.Set("Access-Control-Allow-Origin", allowOrigin)
            headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            headers.Set("Access-Control-Max-Age", "86400")

            // 返回 200 OK
            return ctx.WriteString("")
        }

        // 非 OPTIONS 请求的正常处理
        ctx.Response().Header().Set("Access-Control-Allow-Origin", allowOrigin)
        return ctx.Next()
    }
}
```

**优点**：
- CORS 中间件自动处理 OPTIONS 请求
- 无需手动注册 OPTIONS 路由
- 简单直接

**缺点**：
- 需要修改 CORS 中间件实现

#### 方案 B：调整自动 OPTIONS 注册时机

修改 `router.go` 中 `EnabledAutoOPTIONS` 的注册逻辑，确保在中间件之前：

```go
// 在 RegisterServerFile 或其他注册函数中
func (r *router) RegisterRoute(...) {
    // ...原有逻辑

    // 修改：确保自动 OPTIONS 在中间件之后注册
    // 可能需要调整注册顺序
}
```

**优点**：
- 不需要修改 CORS 中间件
- 保持现有结构

**缺点**：
- 需要深入理解路由注册顺序
- 可能影响其他路由

#### 方案 C：文档化最佳实践

在文档中明确说明 CORS 配置的正确用法：

```markdown
### CORS 配置

**错误用法**（会导致复杂跨域请求失败）：
```go
apiGroup := iweb.HttpServer.Group("/api")
apiGroup.Use(cors.New(""))
iweb.SetEnabledAutoOPTIONS(true)  // ❌ 不要这样做！
```

**正确用法**（手动注册 OPTIONS）：
```go
apiGroup := iweb.HttpServer.Group("/api")
apiGroup.Use(cors.New(""))
apiGroup.OPTIONS("/app", dotweb.DefaultAutoOPTIONSHandler)  // ✅ 推荐

// 或使用支持自动 OPTIONS 的 CORS 中间件（修复后）
apiGroup := iweb.HttpServer.Group("/api")
apiGroup.Use(cors.NewWithAutoOptionsSupport())  // ✅ 推荐
```

### 推荐实施步骤

1. **修复 CORS 中间件**
   - 修改 CORS 中间件，添加 OPTIONS 请求自动处理
   - 测试简单和复杂跨域请求

2. **调整自动 OPTIONS 逻辑**
   - 确保 OPTIONS 处理器在中间件链中正确执行
   - 或者在文档中明确说明最佳实践

3. **添加测试用例**
   - 添加简单跨域请求测试
   - 添加复杂跨域请求测试（带自定义头）
   - 验证 OPTIONS 预检正确

4. **更新文档**
   - 明确 CORS 配置的正确用法
   - 添加示例代码

---

## 总结

### Issue #245（尾斜杠问题）
- **优先级**：高
- **难度**：中等
- **推荐方案**：方案 A（配置选项）
- **预计工作量**：2-3 小时

### Issue #250（CORS 中间件问题）
- **优先级**：中
- **难度**：中等
- **推荐方案**：方案 A（CORS 中间件优化）
- **预计工作量**：2-3 小时

### 后续行动

1. ✅ 分析文档已完成
2. ⏭ 待确认修复方案
3. ⏭ 待实施修复代码
4. ⏭ 待测试验证
5. ⏭ 待创建 PR
