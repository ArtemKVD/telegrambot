package main

//export TELEGRAM_BOT_TOKEN="7327484425:AAFbqP5IYMAIdXoVBzRieM4HESMxEHFtpqY"
import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Panic("TELEGRAM_BOT_TOKEN is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Bot %s is working", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+"BOT")
		reply.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(reply)
		if err != nil {
			log.Println("error send message", err)
		}
	}
}
