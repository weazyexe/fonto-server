package main

import (
	"github.com/weazyexe/fonto-server/internal/app"
)

const configPath = "config/app.yml"

func main() {
	app.Run(configPath)
}
