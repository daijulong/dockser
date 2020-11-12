package bootstrap

import (
	"github.com/daijulong/dockser/commands"
	"github.com/daijulong/dockser/core"
	"github.com/daijulong/dockser/lib"
	"github.com/joho/godotenv"
	"log"
)

func Start()  {
	commands.Run()
}

//加载 .env 文件
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
