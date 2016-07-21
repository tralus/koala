package koala

import (
	"github.com/tralus/koala/knife"
)

type AppModule interface {
	Load(router *knife.Router)
}

func LoadAppModules(router *knife.Router, modules ...AppModule) {
	for _, module:= range modules {
		module.Load(router)
	}
}