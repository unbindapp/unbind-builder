package builder

import "github.com/unbindapp/unbind-builder/config"

type Builder struct {
	config *config.Config
}

func NewBuilder(config *config.Config) *Builder {
	return &Builder{
		config: config,
	}
}
