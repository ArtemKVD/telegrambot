package main

import (
	"log"
	"os"
	"strconv"
	calc "telegrambot/internal/calculate"
	DB "telegrambot/internal/database"
	Redis "telegrambot/internal/redis"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	Redis.InitRedis()
	log.Print("redis connect")
	DB.SetDbConfig()

	defer func() {
		if DB.Db != nil {
			DB.Db.Close()
		}
	}()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	log.Print("Бот запущен")

	type userData struct {
		step    string
		weight  string
		height  string
		gender  string
		program string
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
			gender := update.Message.Text
			if gender != "м" && gender != "ж" {
				sendMessage(bot, chatID, "Ошибка! Введите 'м' или 'ж':")
				continue
			}

			data.gender = gender
			data.step = "program"
			users[userID] = data

			sendMessage(bot, chatID, "choose your program:")

			msg := tgbotapi.NewMessage(chatID, "Выберите программу:")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Похудение"),
					tgbotapi.NewKeyboardButton("Поддержание"),
					tgbotapi.NewKeyboardButton("Набор массы"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("send error: %v", err)
			}
		case "program":
			program := update.Message.Text
			data.program = program
			users[userID] = data

			err := DB.InsertUser(
				username,
				data.weight,
				data.height,
				data.gender,
				data.program,
				calc.Kforlost(data.gender, data.weight, data.height),
				calc.Kforset(data.gender, data.weight, data.height),
				calc.Kforget(data.gender, data.weight, data.height),
			)
			if err != nil {
				log.Printf("user not insert %v", err)
			}
			msg := tgbotapi.NewMessage(chatID, "program set")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			sendMessage(bot, chatID, "your program set")

			delete(users, userID)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("send message error: %v", err)
	}
}
