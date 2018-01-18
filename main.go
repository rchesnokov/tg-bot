package main

import (
	_ "github.com/lib/pq"
	"github.com/rchesnokov/tg-bot/bot"
	"github.com/rchesnokov/tg-bot/service"
	log "github.com/sirupsen/logrus"
)

var (
	env *service.Environment
	db  *service.Database
)

func init() {
	env = service.InitEnviroment()
	db = service.InitDatabase(env.DatabaseURL)

	if env.Mode == "dev" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	bot.New(env.Token)
}
