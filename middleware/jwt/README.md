# jwt
dotweb middleware for jwt.

## 使用：
```
// 设置jwt选项
option := &jwt.Config{
  SigningKey: []byte("devfeel/dotweb"), //must input
  //use cookie
  Extractor: jwt.ExtractorFromCookie,
}
app.AppContext.Set("SimpleJwtConfig", option)
server.Router().GET("/", Index).Use(jwt.NewJWT(option))
```
## 配置：

### Name `string`

jwt token在jwt传输存储中的标识，默认为 `Authorization`

### TTL `time.Duration`

jwt token创建后生存有效期，默认为24小时

### SigningKey `string`

验证token使用的签名字符串

### AuthScheme `string`

用于Authorization header，默认为"Bearer"

### SigningMethod `gojwt.SigningMethod`

加密算法,默认为:jwt.SigningMethodHS256

### ContextKey 

jwt验证通过后，将解密后的token信息中payload的map[string]interface{}存储在dotweb.Context.Items()中

### Extractor `func(name string, ctx dotweb.Context) (string, error)`

提取jwt凭证的方式，默认从header中获取，也可设置为从cookie、querystring中获取，`name`参数是提取token的标识

默认执行 `ExtractorFromHeader`，可设置为为`ExtractorFromCookie` `ExtractorFromQuery`

### ExceptionHandler `dotweb.ExceptionHandle`

验证token过程中出现错误执行的操作， 如用户不设置则默认访问返回401未授权

默认执行 `defaultOnException`

### AddonValidator `func(config *Config, ctx dotweb.Context) error`

附加的token检查器，将在标准token检查通过后执行

### localValidator `func(config *Config, ctx dotweb.Context) error`

本地token检查器

默认执行 `defaultCheckJWT`

建议检查流程如下：

1. 检查是否传递了token
2. 检查token是否可以正常解密
3. 检查token是否过期
4. 检查通过，保存customValue到context中
5. 执行附加检查
