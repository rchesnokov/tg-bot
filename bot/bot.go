package main

import (
	"os"
	"regexp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	environment = "prod"
	token       string
)

func init() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	token = os.Getenv("TOKEN")
	log.Info(token)
	if token == "" {
		log.Fatal("Token wasn't found in environment variable TOKEN")
	}

	r, _ := regexp.Compile("(prod|dev)")
	if r.FindString(os.Getenv("ENV")) != "" {
		environment = os.Getenv("ENV")
	}

	if environment == "dev" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.WithField("account", bot.Self.UserName).Info("Authorized successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.WithFields(log.Fields{
			"username": update.Message.From.UserName,
			"message":  update.Message.Text,
		}).Debug("Incoming message")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
