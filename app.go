package koala

import (
)

// AppModule interface exposes methods to load individual modules.
type AppModule interface {
	Up()
}

// It loads modules that implement AppModule interface
func UpAppModules(modules ...AppModule) {
	for _, module:= range modules {
		module.Up()
	}
}