package main

//go:generate ./scripts/generate.sh

import (
	"github.com/micro-community/micro/v3/cmd"

	_ "github.com/micro-community/micro/v3/cmd/usage"
	// load packages so they can register commands
	_ "github.com/micro-community/micro/v3/client/cli"
	_ "github.com/micro-community/micro/v3/cmd/server"
	_ "github.com/micro-community/micro/v3/cmd/service"
	_ "github.com/micro-community/micro/v3/cmd/usage"
)

func main() {
	cmd.Run()
}
