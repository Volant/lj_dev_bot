package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"

	"./rating"
)

type Configuration struct {
	Token  string
	DbFile string `json:"db_file"`
}

func main() {

	file, _ := os.Open("bot-config.json")
	decoder := json.NewDecoder(file)

	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer file.Close()

	bot, err := tgbotapi.NewBotAPI(configuration.Token)
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		log.Printf("[%s] wrote: [%s]", update.Message.From.UserName, update.Message.Text)

		matched, _ := regexp.MatchString("^[/\\+]", update.Message.Text)
		if matched == true {

			var msg tgbotapi.MessageConfig
			if update.Message.ReplyToMessage != nil {

				rating, err := rating.GetRating(configuration.DbFile, update.Message.ReplyToMessage.From.UserName)

				if err != nil {
					msgText := fmt.Sprintf("Whoops: %s", err)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
					bot.Send(msg)
					continue
				}

				msgText := fmt.Sprintf("%s's rating now %d", update.Message.ReplyToMessage.From.UserName, rating)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "You can't change youself rating")
			}

			bot.Send(msg)
		}
	}
}
