package backend

import (
	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sgorblex/MarkovGenBot/markov"
)

func processMessage(m markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	m.Train(update.Message.Text)
	if update.Message.ReplyToMessage != nil && *update.Message.ReplyToMessage.From == Bot.Self {
		msg := tba.NewMessage(int64(cID), m.Generate())
		msg.ReplyToMessageID = update.Message.MessageID
		Bot.Send(msg)
	}
}
func processCommand(m markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	if update.Message.Command() == "generate" && !m.IsEmpty() {
		msg := tba.NewMessage(int64(cID), m.Generate())
		Bot.Send(msg)
	}
}

func ProcessUpdate(markovs map[ChatID]markov.Markov, update tba.Update) {
	cID := ChatID(update.Message.Chat.ID)
	if update.Message == nil {
		return
	}
	_, exists := markovs[cID]
	if !exists {
		markovs[cID] = make(markov.Markov)
	}
	m := markovs[cID]
	if !update.Message.IsCommand() {
		processMessage(m, update)
	} else {
		processCommand(m, update)
	}
}
