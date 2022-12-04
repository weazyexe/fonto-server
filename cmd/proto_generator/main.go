package main

import (
	"github.com/weazyexe/fonto-server/pkg/proto_generator"
)

const configPath = "config/proto_generator.yml"

func main() {
	proto_generator.Run(configPath)
}
