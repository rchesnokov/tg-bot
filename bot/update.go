package bot

import (
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rchesnokov/tg-bot/horoscope"
	"github.com/rchesnokov/tg-bot/users"
	log "github.com/sirupsen/logrus"
)

// UpdateHandler ... holds link to db and last update from channel
type UpdateHandler struct {
	bot    *tgbotapi.BotAPI
	update tgbotapi.Update
	state  *userState
}

// Process ... processes incoming update
func (uh *UpdateHandler) Process() {
	bot := uh.bot
	state := *uh.state
	update := uh.update

	user := users.GetUser(update.Message.From.UserName)
	text := strings.Replace(update.Message.Text, "@karoshi_bot", "", -1)

	var name string
	if update.Message.From.FirstName != "" {
		name = update.Message.From.FirstName
	} else {
		name = update.Message.From.UserName
	}

	log.WithField("userState", state[user.Name]).Debugf("Current state of user %s", user.Name)

	var msg tgbotapi.MessageConfig

	switch text {
	case "/birthday":
		state[user.Name] = "birthday"
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, name+", введи дату своего рождения в формате dd-mm-yyyy")

	case "/horo":
		birthdate := user.Birthdate
		if !birthdate.Valid {
			state[user.Name] = "birthday+horoscope"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ой, "+name+", я не знаю дату твоего рождения 😥 \nВведи ее в формате dd-mm-yyyy")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, horoscope.Provide(birthdate.String))
		}

	default:
		log.WithField("state", state[user.Name]).Debugf("State of user %s", user.Name)
		switch state[user.Name] {
		case "birthday+horoscope":
			err := processBirthday(user, text)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Дата в неверном формате, попробуй ввести команду заново!")
				msg.ReplyToMessageID = update.Message.MessageID
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Окей, я запомнил! Вот твой гороскоп. \n\n"+horoscope.Provide(user.Birthdate.String))
			}

		case "birthday":
			err := processBirthday(user, text)
			if err != nil {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Дата в неверном формате, попробуй ввести команду заново!")
				msg.ReplyToMessageID = update.Message.MessageID
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Окей, я запомнил!")
				msg.ReplyToMessageID = update.Message.MessageID
			}

		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
		}

		state[user.Name] = ""
	}

	bot.Send(msg)
}

func processBirthday(user *users.User, text string) error {
	date, err := time.Parse("02-01-2006", text)
	if err != nil {
		return err
	}

	user.SetBirthdate(date.Format("2006-01-02"))

	return nil
}
