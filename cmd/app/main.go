package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	DB "telegrambot/internal/database"
	"telegrambot/internal/limits"
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

	log.Print("bot starting")

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
				sendMessage(bot, chatID, "insert м or ж")
				continue
			}

			data.gender = gender
			data.step = "program"
			users[userID] = data

			sendMessage(bot, chatID, "choose your program:")

			msg := tgbotapi.NewMessage(chatID, "choose program:")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("lost"),
					tgbotapi.NewKeyboardButton("set"),
					tgbotapi.NewKeyboardButton("get"),
				),
			)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("send error: %v", err)
			}
		case "program":
			program := update.Message.Text
			data.program = program
			users[userID] = data

			dailyLimits, err := limits.Calculate(
				data.gender,
				data.weight,
				data.height,
				program,
			)
			if err != nil {
				log.Printf("limits error: %v", err)
				sendMessage(bot, chatID, "limits not calculated")
				continue
			}

			if err := Redis.SetUserLimits(username, dailyLimits); err != nil {
				log.Printf("redis error: %v", err)
				sendMessage(bot, chatID, "limits not set")
				continue
			}

			msg := fmt.Sprintf(`
				PROGRAM %v
				DAILY LIMITS:
				calories: %v 
				proteins: %v 
				fats: %v 
				carbs: %v 
				`,
				program,
				dailyLimits.Calories,
				dailyLimits.Proteins,
				dailyLimits.Fats,
				dailyLimits.Carbs)

			msgConfig := tgbotapi.NewMessage(chatID, msg)
			msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msgConfig); err != nil {
				log.Printf("send error: %v", err)
			}

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
