package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/devfeel/dotweb"
	gojwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

var (
	SigningMethodHS256 *gojwt.SigningMethodHMAC
	SigningMethodHS384 *gojwt.SigningMethodHMAC
	SigningMethodHS512 *gojwt.SigningMethodHMAC

	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = &Config{
		TTL:              time.Hour * 24,
		SigningMethod:    SigningMethodHS256,
		ContextKey:       "dotjwt-user",
		Name:             dotweb.HeaderAuthorization,
		AuthScheme:       "Bearer",
		Extractor:        ExtractorFromHeader,
		ExceptionHandler: defaultOnException,
		validator:        defaultCheckJWT,
	}
)

func init() {
	SigningMethodHS256 = gojwt.SigningMethodHS256
	SigningMethodHS384 = gojwt.SigningMethodHS384
	SigningMethodHS512 = gojwt.SigningMethodHS512
}

type (
	// TokenExtractor extract jwt token
	tokenExtractor func(name string, ctx dotweb.Context) (string, error)

	// addonValidator addon validator, after jwt standard validator pass
	addonValidator func(config *Config, ctx dotweb.Context) error

	// localValidator jwt standard validator
	localValidator func(config *Config, ctx dotweb.Context) error

	//Config JWT Middleware config
	Config struct {
		// TTL time to live, int64 nanosecond count
		TTL time.Duration
		// Name token name, default: Authorization
		Name string
		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string
		// Signing key to validate token
		SigningKey interface{}
		// ExceptionHandler validate error handler, default: defaultOnError
		ExceptionHandler dotweb.ExceptionHandle

		// Extractor extract jwt token, default extract from header: defaultExtractorFromHeader
		Extractor tokenExtractor
		// EnableAuthOnOptions http option method validate switch
		EnableAuthOnOptions bool
		// SigningMethod sign method, default: HS256
		SigningMethod gojwt.SigningMethod
		// ContextKey Context key to store user information from the token into context.
		// Optional. Default value "dotjwt-user".
		ContextKey string
		//AddonValidator addon validator will handle after standard validator
		AddonValidator addonValidator
		// validator custom validator suggestion flow：
		// 1. extract token string
		// 2. check token sign
		// 3. check token ttl
		// 4. save custom value to conext after check passed
		// 5. handle addon validator
		validator localValidator
		// validationKeyGetter
		keyFunc gojwt.Keyfunc
	}
)

//jwt中间件
type JWTMiddleware struct {
	dotweb.BaseMiddlware
	config *Config
}

func (m *JWTMiddleware) Handle(ctx dotweb.Context) error {
	if err := m.config.validator(m.config, ctx); err != nil {
		m.config.ExceptionHandler(ctx, err)
		return nil
	}
	m.Next(ctx)
	return nil
}

// New create a JWT Middleware
func NewJWT(config *Config) *JWTMiddleware {
	if config.Name == "" {
		config.Name = DefaultJWTConfig.Name
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}
	if config.ContextKey == "" {
		config.ContextKey = config.Name
	}
	if config.Extractor == nil {
		config.Extractor = DefaultJWTConfig.Extractor
	}
	if config.SigningMethod == nil {
		config.SigningMethod = DefaultJWTConfig.SigningMethod
	}

	if config.ExceptionHandler == nil {
		config.ExceptionHandler = DefaultJWTConfig.ExceptionHandler
	}
	if config.validator == nil {
		config.validator = DefaultJWTConfig.validator
	}

	if config.SigningKey == "" {
		panic("jwt middleware requires signing key")
	} else {
		config.keyFunc = func(t *gojwt.Token) (interface{}, error) {
			// Check the signing method
			if t.Method.Alg() != config.SigningMethod.Alg() {
				return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
			}
			return config.SigningKey, nil
		}
	}

	return &JWTMiddleware{config: config}
}

// defaultOnException default exception handler
// return false will break conext, return true will handle next context
func defaultOnException(ctx dotweb.Context, err error) {
	ctx.WriteStringC(http.StatusUnauthorized, err.Error())
}

// defaultCheckJWT execlude token check flow, or returns error
// 1. extract token string
// 2. check token sign
// 3. check token ttl
// 4. save custom value to conext after check passed
// 5. handle addon validator
func defaultCheckJWT(config *Config, ctx dotweb.Context) error {
	req := ctx.Request()

	if !config.EnableAuthOnOptions {
		if req.Method == "OPTIONS" {
			return nil
		}
	}
	// extract token
	token, err := config.Extractor(config.Name, ctx)

	if err != nil {
		return fmt.Errorf("Error extracting token: %v", err)
	}
	if token == "" {
		// no token
		return fmt.Errorf("Required authorization token not found")
	}

	// parse token value
	parsedToken, err := gojwt.Parse(token, config.keyFunc)
	if err != nil {
		return fmt.Errorf("Error parsing token: %v", err)
	}

	if config.SigningMethod != nil && config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		return fmt.Errorf("Error validating token algorithm: %s", message)
	}

	if !parsedToken.Valid {
		return fmt.Errorf("Token is invalid")
	}
	// save custom value to context
	claims := parsedToken.Claims.(gojwt.MapClaims)
	ctx.Items().Set(config.ContextKey, map[string]interface{}(claims))

	// handle addon validator
	if config.AddonValidator != nil {
		err := config.AddonValidator(config, ctx)
		if err != nil {
			return fmt.Errorf("JWT Addon validate check failed: %s", err)
		}
	}

	return nil
}

// GeneratorToken generate token by custom value and token ttl
// 标准中注册的声明 (建议但不强制使用) ：
// iss: jwt签发者
// sub: jwt所面向的用户
// aud: 接收jwt的一方
// exp: jwt的过期时间，这个过期时间必须要大于签发时间
// nbf: 定义在什么时间之前，该jwt都是不可用的.
// iat: jwt的签发时间
// jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
// 我们这里默认使用exp控制有效期
func GeneratorToken(config *Config, payload map[string]interface{}) (string, error) {
	claims := gojwt.MapClaims(payload)
	//特别的，如果未设置iat与exp，设置默认值
	//设置iat
	if _, isok := claims["iat"]; !isok {
		claims["iat"] = time.Now().Unix()
	}
	//设置exp
	if _, isok := claims["exp"]; !isok {
		claims["exp"] = time.Now().Add(config.TTL).Unix()
	}

	token := gojwt.NewWithClaims(config.SigningMethod, claims)
	// sign token and get the complete encoded token as a string
	return token.SignedString(config.SigningKey)
}

// ExtractorFromHeader extract token from header
// if use header mode, please add  "bearer "
func ExtractorFromHeader(name string, ctx dotweb.Context) (string, error) {
	authHeader := ctx.Request().Header.Get(name)
	if authHeader == "" || len(authHeader) <= 7 {
		return "", nil
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// ExtractorFromCookie extract token from cookie
func ExtractorFromCookie(name string, ctx dotweb.Context) (string, error) {
	authCookie, err := ctx.ReadCookie(name)
	if err != nil || authCookie == nil || authCookie.Value == "" {
		return "", nil
	}
	return authCookie.Value, nil
}

// ExtractorFromQuery extract token from query string
func ExtractorFromQuery(name string, ctx dotweb.Context) (string, error) {
	token := ctx.QueryString(name)
	if token == "" {
		return "", nil
	}
	return token, nil
}
