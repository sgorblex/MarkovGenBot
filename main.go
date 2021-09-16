package main

// TODO
// actual commands
// function for generating text of x words. When a sentence terminates (next = ""), add . and start again (prev = "")
// markov should be populated per chat
// only what is convenient to keep in memory should be kept, the rest goes on file/db

import (
	"bufio"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	markov "github.com/SimpoLab/SimpoBot/markov"
)

func getApiKey() string {
	keyFile, err := os.Open("api_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer keyFile.Close()

	sc := bufio.NewScanner(keyFile)
	if !sc.Scan() {
		log.Fatal("ERROR: API Key file is empty.")
	}
	return sc.Text()
}
func main() {
	bot, err := tgbotapi.NewBotAPI(getApiKey())
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	u.Timeout = 2
	updates, err := bot.GetUpdatesChan(u)

	fmt.Fprintln(os.Stderr, "bot started")
	m := make(markov.Markov)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// this if condition is only for testing purposes and will be changed
		if !update.Message.IsCommand() {
			m.Train(update.Message.Text)
		} else {
			if !m.IsEmpty() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, m.Generate())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}

	}
}
