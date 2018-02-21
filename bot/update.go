package bot

import (
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rchesnokov/tg-bot/features"
	"github.com/rchesnokov/tg-bot/users"
	log "github.com/sirupsen/logrus"
)

// UpdateHandler ... holds users' state and last update from channel
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

	chatID := update.Message.Chat.ID
	messageID := update.Message.MessageID
	text := strings.Replace(update.Message.Text, "@karoshi_bot", "", -1)

	username := update.Message.From.UserName
	firstname := update.Message.From.FirstName

	createMessage := createMessage(chatID)

	user := users.FindByUsername(username)
	if user == nil {
		user = users.Create(username, firstname)
	}

	name := user.GetName()

	var msg tgbotapi.MessageConfig

	log.WithField("userState", state[user.Username]).Debugf("Current state of user %s", user.Username)

	switch text {
	case "/birthday":
		state[user.Username] = "birthday"
		msg = createMessage(name + ", введи дату своего рождения в формате dd-mm-yyyy")

	case "/bydlo":
		msg = createMessage(features.PrintSwearingRating())
		msg.ParseMode = "HTML"

	case "/horo":
		birthdate := user.Birthdate
		if birthdate == "" {
			state[user.Username] = "birthday+horoscope"
			msg = createMessage("Ой, " + name + ", я не знаю дату твоего рождения 😥 \nВведи ее в формате dd-mm-yyyy")
		} else {
			msg = createMessage(features.ProvideHoroscope(birthdate))
		}

	default:
		log.WithField("state", state[user.Username]).Debugf("State of user %s", user.Username)

		switch state[user.Username] {
		case "birthday+horoscope":
			err := handleBirthday(user, text)
			if err != nil {
				msg = createMessage("Дата в неверном формате, попробуй ввести команду заново!")
				msg.ReplyToMessageID = messageID
			} else {
				msg = createMessage("Окей, я запомнил! Вот твой гороскоп. \n\n" + features.ProvideHoroscope(user.Birthdate))
			}

		case "birthday":
			err := handleBirthday(user, text)
			if err != nil {
				msg = createMessage("Дата в неверном формате, попробуй ввести команду заново!")
				msg.ReplyToMessageID = messageID
			} else {
				msg = createMessage("Окей, я запомнил!")
				msg.ReplyToMessageID = messageID
			}

		default:
			handleSwearing(user, text)
		}

		state[user.Username] = ""
	}

	bot.Send(msg)
}

func createMessage(chatID int64) func(string) tgbotapi.MessageConfig {
	return func(message string) tgbotapi.MessageConfig {
		return tgbotapi.NewMessage(chatID, message)
	}
}

func handleBirthday(user *users.User, text string) error {
	date, err := time.Parse("02-01-2006", text)
	if err != nil {
		return err
	}

	user.SetBirthdate(date.Format("2006-01-02"))

	return nil
}

func handleSwearing(user *users.User, text string) {
	c := features.FilterSwearing(text)
	if c > 0 {
		user.SetSwearing(c)
	}
}
