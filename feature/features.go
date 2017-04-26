package feature

type Feature struct {
	CROSConfig *CROSConfig
}

func NewFeature() *Feature {
	return &Feature{
		CROSConfig: NewCORSConfig(),
	}
}

//set Enabled CROS true, with default config
func (f *Feature) SetEnabledCROS() *CROSConfig {
	if f.CROSConfig == nil {
		f.CROSConfig = NewCORSConfig()
	}
	f.CROSConfig.EnabledCROS = true
	f.CROSConfig.UseDefault()
	return f.CROSConfig
}

//set Disabled CROS false
func (f *Feature) SetDisabledCROS() {
	if f.CROSConfig == nil {
		f.CROSConfig = NewCORSConfig()
	}
	f.CROSConfig.EnabledCROS = false
}
