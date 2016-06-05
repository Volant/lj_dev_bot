package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Volant/lj_dev_bot/ljtop"
	"github.com/Volant/lj_dev_bot/rating"
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

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// if update.InlineQuery == nil { // if no inline query, ignore it
		// 	continue
		// }

		log.Printf("[%s] wrote: [%s]", update.Message.From.UserName, update.Message.Text)

		req_rating := regexp.MustCompile("^/?\\+")
		req_ljtop := regexp.MustCompile("^/ljtop")
		req_help := regexp.MustCompile("^/help")
		req_test_feature := regexp.MustCompile("^test")
		switch {
		case req_test_feature.MatchString(update.Message.Text):
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "A")

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("TEST", "data"), tgbotapi.NewInlineKeyboardButtonData("TEST_DATA", "test_data")))

			// msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
			// tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("TEST KEYBOARD")))
			// tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("TEST", "data")))

			// resp, err := bot.AnswerCallbackQuery(msg)
			// if err != nil {
			// 	fmt.Println(err)
			// } else {
			// 	fmt.Println(resp)
			// }
			bot.Send(msg)
		case req_help.MatchString(update.Message.Text):
			var msg tgbotapi.MessageConfig
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, `It is LJ Tiny SWD's bot. Available commands
/+ in reply to message - increase rating
/ljtop [COUNTRY] - get LJ Rating for selected COUNTRY
/help - this text`)

			bot.Send(msg)
		case req_ljtop.MatchString(update.Message.Text):
			var msg tgbotapi.MessageConfig
			rating := ljtop.GetLJTop("cyr")
			for position, rating_entry := range rating.Result.Rating {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d. %s", position+1, rating_entry.PostUrl))
				bot.Send(msg)
			}
		case req_rating.MatchString(update.Message.Text):
			var msg tgbotapi.MessageConfig
			if update.Message.ReplyToMessage != nil {

				rating, err := rating.GetRating(configuration.DbFile, update.Message.ReplyToMessage.From.UserName)

				if err != nil {
					msgText := fmt.Sprintf("Whoops: %s", err)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
					bot.Send(msg)
					continue
				}

				msgText := fmt.Sprintf("%s's rating is now %d", update.Message.ReplyToMessage.From.UserName, rating)
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "You can't change youself rating")
			}

			bot.Send(msg)
		}
	}
}
