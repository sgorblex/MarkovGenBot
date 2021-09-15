package main

// TODO:
// !!! when we generated a word with no Followings and we have to generate more, (finish the sentence with . and) start again with ""
// words that terminate a text could have "" as Following, and when that happens markov generates a . if it has to generate more. Or, even better, instead of generating n words it could generate a message, which ends when "" is next. Or might implement both functions.
// markov should be populated per chat
// only what is convenient to keep in memory should be kept, the rest goes on file/db
// markov should be in its own package
// functions that take a Folloing should instead be methods of Markov and take a string, which will then define the Following
// check that it does not train itself

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Following map[string]uint
type Markov map[string]Following

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

func (m Markov) IsEmpty() bool {
	return len(m) == 0
}

func (m Markov) Populate(text string) {
	prev := ""
	for _, word := range strings.Fields(text) {
		if m[prev] == nil {
			m[prev] = make(Following)
		}
		m[prev][word]++
		prev = word
	}
}

func sumOfProb(pref Following) uint {
	var i uint
	for _, prob := range pref {
		i += prob
	}
	return i
}

func (m Markov) genWord(pref Following) string {
	var i, r uint
	// to be fixed
	sop := sumOfProb(pref)
	if sop == 0 {
		return "Failed. Fix this bug, bro."
	}
	r = uint(rand.Uint64()) % sop
	for suff, prob := range pref {
		i += prob
		if i >= r {
			return suff
		}
	}
	log.Panic("assertion error")
	return ""
}

func (m Markov) Generate(nWords int) string {
	var (
		res  string
		r, j int
		prev string
	)
	r = rand.Intn(len(m))
	for prev = range m {
		if r >= j {
			break
		}
		j++
	}
	for i := 0; i < nWords; i++ {
		word := m.genWord(m[prev])
		res += word + " "
		prev = word
	}
	return res
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

	m := make(Markov)
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
