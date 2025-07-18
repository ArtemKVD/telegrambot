package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	calc "telegrambot/internal/calculate"
	DB "telegrambot/internal/database"
	limits "telegrambot/internal/limits"
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
			sendMessage(bot, chatID, "This bot will help you monitor your diet, now you need to enter some data so that it calculates the norm of calories, proteins, fats and carbohydrates for losing weight, maintaining and gaining weight, after which you can add meals and the bot will show you how much you have left to eat so as not to go beyond the norm")
			sendMessage(bot, chatID, "use commad /add_meal for add meal or /last after last meal of the day")
			sendMessage(bot, chatID, "Insert you kilo weight")
			continue
		}
		if update.Message.IsCommand() && update.Message.Command() == "addmeal" {
			users[userID] = userData{step: "meal_input"}
			sendMessage(bot, chatID, "Insert CPFC in format 0 0 0 0")
			continue
		}
		if update.Message.IsCommand() && update.Message.Command() == "last" {
			userData, err := DB.GetUserData(DB.Db, username)
			log.Println(userData.Gender, userData.Height, userData.Program, userData.Weight)
			if err != nil {
				log.Printf("get data error: %v", err)
				sendMessage(bot, chatID, "get data error write /start")
				continue
			}

			dailyLimits, err := limits.Calculate(
				userData.Gender,
				userData.Weight,
				userData.Height,
				userData.Program,
			)
			log.Println(userData.Gender, userData.Weight, userData.Height, userData.Program)
			if err != nil {
				log.Printf("limits error %v", err)
				sendMessage(bot, chatID, "limits error")
				continue
			}

			if err := Redis.SetUserLimits(username, dailyLimits); err != nil {
				log.Printf("Redis error %v", err)
				sendMessage(bot, chatID, "update limits error")
				continue
			}

			sendMessage(bot, chatID, fmt.Sprintf(`
				limits updated
				calories: %d 
				proteins: %d 
				fats: %d 
				carbs: %d `,
				dailyLimits.Calories,
				dailyLimits.Proteins,
				dailyLimits.Fats,
				dailyLimits.Carbs))

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
			if gender != "m" && gender != "w" {
				sendMessage(bot, chatID, "insert m or w")
				continue
			}

			data.gender = gender
			data.step = "program"
			users[userID] = data

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
			data.step = "meal_input"

			dailyLimits, err := limits.Calculate(
				data.gender,
				data.weight,
				data.height,
				program,
			)
			log.Println(data.gender, data.height, data.program, data.weight)
			DB.InsertUser(username, data.weight, data.height, data.gender, data.program, calc.Kforset(data.gender, data.height, data.weight), calc.Kforget(data.gender, data.height, data.weight), calc.Kforlost(data.gender, data.height, data.weight))
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
		case "meal_input":
			parts := strings.Fields(update.Message.Text)
			if len(parts) != 4 {
				sendMessage(bot, chatID, "Insert 4 numbers")
				continue
			}

			calories, err1 := strconv.Atoi(parts[0])
			proteins, err2 := strconv.Atoi(parts[1])
			fats, err3 := strconv.Atoi(parts[2])
			carbs, err4 := strconv.Atoi(parts[3])

			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				sendMessage(bot, chatID, "error")
				continue
			}

			remaining, err := Redis.SubtractMeal(username, calories, proteins, fats, carbs)
			if err != nil {
				log.Printf("update error %v", err)
				sendMessage(bot, chatID, "error data update")
				continue
			}

			response := fmt.Sprintf(`
					meal add:
					calories: %d
					proteins: %d 
					fats: %d 
					carbs: %d 
					`,
				remaining.Calories,
				remaining.Proteins,
				remaining.Fats,
				remaining.Carbs,
			)

			sendMessage(bot, chatID, response)
			delete(users, userID)
			continue
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("send message error: %v", err)
	}
}
