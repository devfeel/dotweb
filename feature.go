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
func (f *xFeatureTools) SetCROSConfig(ctx *HttpContext, c *feature.CROSConfig) {
	ctx.SetHeader(HeaderAccessControlAllowOrigin, c.AllowedOrigins)
	ctx.SetHeader(HeaderAccessControlAllowMethods, c.AllowedMethods)
	ctx.SetHeader(HeaderAccessControlAllowHeaders, c.AllowedHeaders)
	ctx.SetHeader(HeaderAccessControlAllowCredentials, strconv.FormatBool(c.AllowCredentials))
	ctx.SetHeader(HeaderP3P, c.AllowedP3P)
}
