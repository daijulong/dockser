package main

import (
	"github.com/daijulong/dockser/v2/bootstrap"
)

func main() {
	bootstrap.LoadEnv()
	bootstrap.Start()
}
