package main

import (
	"os"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/servers"
	"github.com/codepnw/ecommerce/pkg/database"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())
	// fmt.Println(cfg.Db())

	db := database.DbConnect(cfg.Db())
	defer db.Close()
	
	servers.NewServer(cfg, db).Start()
}