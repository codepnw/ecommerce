package mytests

import (
	"encoding/json"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/servers"
	"github.com/codepnw/ecommerce/pkg/database"
)

func SetupTest() servers.IModuleFactory {
	cfg := config.LoadConfig("../.env.test")

	db := database.DbConnect(cfg.Db())
	defer db.Close()

	s := servers.NewServer(cfg, db)
	return servers.InitModule(nil, s.GetServer(), nil)
}

func CompressToJson(obj any) string {
	result, _ := json.Marshal(&obj)
	return string(result)
}
