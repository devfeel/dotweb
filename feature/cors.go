package feature

//CROS配置
type CROSConfig struct {
	EnabledCROS      bool
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	AllowCredentials bool
	AllowedP3P       string
}

func NewCORSConfig() *CROSConfig {
	return &CROSConfig{}
}

func (c *CROSConfig) UseDefault() *CROSConfig {
	c.AllowedOrigins = "*"
	c.AllowedMethods = "GET, POST, PUT, DELETE, OPTIONS"
	c.AllowedHeaders = "Content-Type"
	c.AllowedP3P = "CP=\"CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA PRE COM NAV OTC NOI DSP COR\""
	return c
}

func (c *CROSConfig) SetOrigin(origins string) *CROSConfig {
	c.AllowedOrigins = origins
	return c
}

func (c *CROSConfig) SetMethod(methods string) *CROSConfig {
	c.AllowedMethods = methods
	return c
}

func (c *CROSConfig) SetHeader(headers string) *CROSConfig {
	c.AllowedHeaders = headers
	return c
}
func (c *CROSConfig) SetAllowCredentials(flag bool) *CROSConfig {
	c.AllowCredentials = flag
	return c
}
