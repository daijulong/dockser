package main

import (
	"github.com/daijulong/dockser/bootstrap"
)

func main() {
	bootstrap.LoadEnv()
	bootstrap.Start()
}
