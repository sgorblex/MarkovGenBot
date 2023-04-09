package backend

import (
	"bufio"
	"log"
	"os"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	apiKeyFileName string = "api_key.txt"
)

type ChatID int64

var Bot *tba.BotAPI

func init() {
	var err error
	Bot, err = tba.NewBotAPI(getApiKey())
	if err != nil {
		log.Panic(err)
	}
}

func getApiKey() string {
	keyFile, err := os.Open(apiKeyFileName)
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

func GetUpdates() tba.UpdatesChannel {
	u := tba.NewUpdate(0)
	u.Timeout = 60
	return Bot.GetUpdatesChan(u)
}
