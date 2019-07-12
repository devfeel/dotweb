package dotweb

// Plugin a interface for app's global plugin
type Plugin interface {
	Name() string
	Run() error
	IsValidate() bool
}
