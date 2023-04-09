package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sgorblex/MarkovGenBot/markov"
)

const (
	baseDataPath   string = "data"
	whitelistsPath string = "whitelist.json"
)

var whitelist []ChatID
var useWhitelist bool

func init() {
	file, err := os.Open(whitelistsPath)
	if err != nil {
		log.Printf("Not using whitelist, reason: %v.\n", err)
		return
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(raw, &whitelist)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Parsed whitelist: %v.\n", whitelist)
	useWhitelist = true
}

type Tables map[ChatID]markov.Markov

func processMessage(m markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	log.Printf("Training for chat %v.\n", cID)
	m.Train(update.Message.Text)
	if update.Message.ReplyToMessage != nil && *update.Message.ReplyToMessage.From == Bot.Self {
		log.Printf("Generating for chat %v.\n", cID)
		msg := tba.NewMessage(int64(cID), m.Generate())
		msg.ReplyToMessageID = update.Message.MessageID
		Bot.Send(msg)
	}
}
func processCommand(m markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	if update.Message.Command() == "generate" && !m.IsEmpty() {
		log.Printf("Generating for chat %v.\n", cID)
		msg := tba.NewMessage(int64(cID), m.Generate())
		Bot.Send(msg)
	}
}

func (t Tables) fetchOrCreate(cID ChatID) (markov.Markov, error) {
	m, exists := t[cID]
	if exists {
		return m, nil
	}
	filePath := baseDataPath + "/" + strconv.Itoa(int(cID)) + ".json"
	if _, err := os.Stat(filePath); err != nil {
		log.Printf("No memory for chat %v; creating new.\n", cID)
		t[cID] = make(markov.Markov)
		return t[cID], nil
	}
	log.Printf("Found persistent memory for chat %v; unmarshaling.\n", cID)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}
	t[cID] = m
	return m, nil
}

func isWhitelisted(cID ChatID) bool {
	for _, candidate := range whitelist {
		if cID == candidate {
			return true
		}
	}
	return false
}

func ProcessUpdate(markovs Tables, update tba.Update) error {
	if update.Message == nil {
		return nil
	}
	cID := ChatID(update.Message.Chat.ID)
	if useWhitelist && !isWhitelisted(cID) {
		return fmt.Errorf("Skipping update from chat %v as it's now whitelisted", cID)
	}
	m, err := markovs.fetchOrCreate(cID)
	if err != nil {
		return fmt.Errorf("Error with update on chat %v: %v.", cID, err)
	}
	if !update.Message.IsCommand() {
		processMessage(m, update)
	} else {
		processCommand(m, update)
	}
	return nil
}
