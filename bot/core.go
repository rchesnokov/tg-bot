package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type userState map[string]string

// New ... creates new bot
func New(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.WithField("account", bot.Self.UserName).Info("Authorized successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	us := make(userState)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.WithFields(log.Fields{
			"username": update.Message.From.UserName,
			"message":  update.Message.Text,
		}).Debug("Incoming message")

		upd := &UpdateHandler{
			bot:    bot,
			state:  &us,
			update: update,
		}

		upd.Process()
	}
}
