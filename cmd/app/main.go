package main

import (
	"log"
	"os"
	"strconv"
	DB "telegrambot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	DB.SetDbConfig()
	defer DB.Db.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Print("Бот запущен")

	type userData struct {
		step   string
		weight string
		height string
		gender string
	}
	users := make(map[int64]userData)

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		username := update.Message.From.UserName
		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			users[userID] = userData{step: "weight"}
			sendMessage(bot, chatID, "Insert you kilo weight")
			continue
		}

		data, exists := users[userID]
		if !exists {
			continue
		}

		switch data.step {
		case "weight":
			_, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				sendMessage(bot, chatID, "error")
				continue
			}
			data.weight = update.Message.Text
			data.step = "height"
			users[userID] = data
			sendMessage(bot, chatID, "Insert your height")

		case "height":
			_, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				sendMessage(bot, chatID, "error")
				continue
			}
			data.height = update.Message.Text
			data.step = "gender"
			users[userID] = data
			sendMessage(bot, chatID, "set your gender:")

		case "gender":
			data.gender = update.Message.Text
			users[userID] = data

			err := DB.InsertUser(
				username,
				data.weight,
				data.height,
				data.gender,
				1, 1, 1,
			)

			if err != nil {
				sendMessage(bot, chatID, "data not saved")
				log.Printf("DB error: %v", err)
			} else {
				sendMessage(bot, chatID, "Data saved")
			}
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}
