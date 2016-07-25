package koala

import (
)

type AppModule interface {
	Up()
}

func UpAppModules(modules ...AppModule) {
	for _, module:= range modules {
		module.Up()
	}
}