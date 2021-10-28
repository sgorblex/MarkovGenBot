package main

// TODO
// function for generating text of x words. When a sentence terminates (next = ""), add . and start again (prev = "")
// only what is convenient to keep in memory should be kept, the rest goes on file/db
// list of allowed chats
// logging (with chat codes)
// cache of sum of probs
// write to file every x minutes (protection from system crashes)
// use database instead of json, or at least use a different file per chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	markov "github.com/SimpoLab/MarkovBot/markov"
)

const (
	dataFileName   string = "data.json"
	apiKeyFileName string = "api_key.txt"
)

type ChatID int64

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

func gracefulExit(m map[ChatID]markov.Markov) {
	fmt.Fprintf(os.Stderr, "\nexiting gracefully...\n")
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(dataFileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	n, err := file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Truncate(int64(n))
}

func markovsFromFile(filename string) (map[ChatID]markov.Markov, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var m map[ChatID]markov.Markov
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
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

	// markovs := make(markov.Markov)
	var markovs map[ChatID]markov.Markov
	if _, err := os.Stat(dataFileName); err == nil {
		markovs, err = markovsFromFile(dataFileName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		markovs = make(map[ChatID]markov.Markov)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		gracefulExit(markovs)
		os.Exit(0)
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// this if condition is only for testing purposes and will be changed
		var m markov.Markov
		if markovs[ChatID(update.Message.Chat.ID)] == nil {
			markovs[ChatID(update.Message.Chat.ID)] = make(markov.Markov)
		}
		m = markovs[ChatID(update.Message.Chat.ID)]
		if !update.Message.IsCommand() {
			m.Train(update.Message.Text)
		} else {
			if update.Message.Command() == "generate" && !m.IsEmpty() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, m.Generate())
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}

	}
}
