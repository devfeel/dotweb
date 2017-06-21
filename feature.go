package dotweb

import (
	"github.com/devfeel/dotweb/feature"
	"strconv"
)

type xFeatureTools struct{}

var FeatureTools *xFeatureTools

func init() {
	FeatureTools = new(xFeatureTools)
}

//set CROS config on HttpContext
func (f *xFeatureTools) SetCROSConfig(ctx Context, c *feature.CROSConfig) {
	ctx.Response().SetHeader(HeaderAccessControlAllowOrigin, c.AllowedOrigins)
	ctx.Response().SetHeader(HeaderAccessControlAllowMethods, c.AllowedMethods)
	ctx.Response().SetHeader(HeaderAccessControlAllowHeaders, c.AllowedHeaders)
	ctx.Response().SetHeader(HeaderAccessControlAllowCredentials, strconv.FormatBool(c.AllowCredentials))
	ctx.Response().SetHeader(HeaderP3P, c.AllowedP3P)
}
