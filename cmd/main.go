package main

import (
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/builder"
)

func main() {
	godotenv.Load()

	builder := builder.NewBuilder(
		config.NewConfig(),
	)

	builder.BuildWithNixpacks()
}
