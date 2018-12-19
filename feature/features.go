package feature

type Feature struct {
	CROSConfig *CROSConfig
}

func NewFeature() *Feature {
	return &Feature{
		CROSConfig: NewCORSConfig(),
	}
}

// SetEnabledCROS enable CROS, with default config
func (f *Feature) SetEnabledCROS() *CROSConfig {
	if f.CROSConfig == nil {
		f.CROSConfig = NewCORSConfig()
	}
	f.CROSConfig.EnabledCROS = true
	f.CROSConfig.UseDefault()
	return f.CROSConfig
}

// SetDisabledCROS disable CROS
func (f *Feature) SetDisabledCROS() {
	if f.CROSConfig == nil {
		f.CROSConfig = NewCORSConfig()
	}
	f.CROSConfig.EnabledCROS = false
}
