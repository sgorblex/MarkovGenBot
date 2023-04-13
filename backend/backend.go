package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sgorblex/MarkovGenBot/markov"
)

const (
	baseDataPath   string        = "data"
	whitelistsPath string        = "whitelist.json"
	oldThreshold   time.Duration = 24 * time.Hour
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

type TimedMarkov struct {
	access time.Time
	mark   markov.Markov
}
type Tables map[ChatID]TimedMarkov

func processMessage(m markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	log.Printf("Training for chat %v.\n", cID)
	m.Train(update.Message.Text)
	if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.From.ID == Bot.Self.ID {
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

func (t Tables) UnloadOld() {
	t.Persist()
	// now is only requested once; might lose information since as of now the map is not locked
	oldTime := time.Now().Add(-oldThreshold)
	for k, v := range t {
		if v.access.Before(oldTime) {
			log.Printf("Table of chat %v unloaded since last access was on %v.\n", k, v.access)
			delete(t, k)
		}
	}
}

func (t Tables) fetchOrCreate(cID ChatID) (markov.Markov, error) {
	tm, exists := t[cID]
	if exists {
		t[cID] = TimedMarkov{time.Now(), tm.mark}
		return tm.mark, nil
	}
	filePath := baseDataPath + "/" + strconv.Itoa(int(cID)) + ".json"
	if _, err := os.Stat(filePath); err != nil {
		log.Printf("No memory for chat %v; creating new.\n", cID)
		m := make(markov.Markov)
		t[cID] = TimedMarkov{time.Now(), m}
		return m, nil
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
	m := make(markov.Markov)
	err = json.Unmarshal(raw, &m)
	if err != nil {
		return nil, err
	}
	t[cID] = TimedMarkov{time.Now(), m}
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
