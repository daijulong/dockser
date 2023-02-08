package bootstrap

import (
	"github.com/daijulong/dockser/commands"
	"github.com/daijulong/dockser/core"
	"github.com/daijulong/dockser/lib"
	"github.com/joho/godotenv"
	"log"
)

// Start 启动
func Start() {
	commands.Run()
}

// LoadEnv 加载 .env 文件
func LoadEnv() {
	core.Envs = make(map[string]string)
	if lib.IsFile(".env") {
		envs, err := godotenv.Read(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		core.Envs = envs
	}
}
