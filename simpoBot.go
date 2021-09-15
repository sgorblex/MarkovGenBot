package main

// TODO
// !!! when we generated a word with no Followings and we have to generate more, (finish the sentence with . and) start again with ""
// words that terminate a text could have "" as Following, and when that happens markov generates a . if it has to generate more. Or, even better, instead of generating n words it could generate a message, which ends when "" is next. Or might implement both functions.
// markov should be populated per chat
// only what is convenient to keep in memory should be kept, the rest goes on file/db
// check that it does not train itself
// see markov package

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	markov "github.com/SimpoLab/SimpoBot/markov" // verify project structure
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
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	m := make(markov.Markov)
	rand.Seed(time.Now().UnixNano())
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// this if condition is only for testing purposes and will be changed
		if update.Message.IsCommand() && !m.IsEmpty() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, m.Generate(200))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		} else {
			m.Populate(update.Message.Text)
		}

	}
}
