package markov

import (
	"math/rand"
	"strings"
	"time"

	assert "github.com/SimpoLab/MarkovBot/assert"
)

type Following map[string]uint
type Markov map[string]Following

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (m Markov) IsEmpty() bool {
	return len(m) == 0
}

func (m Markov) Train(text string) {
	var prev string
	for _, word := range strings.Fields(text) {
		if m[prev] == nil {
			m[prev] = make(Following)
		}
		m[prev][word]++
		prev = strings.ToLower(word)
	}
	if m[prev] == nil {
		m[prev] = make(Following)
	}
	m[prev][""]++
}

func (m Markov) sumOfProb(word string) uint {
	var i uint
	for _, prob := range m[word] {
		i += prob
	}
	return i
}

func (m Markov) genWord(prev string) string {
	lowPrev := strings.ToLower(prev)
	sop := m.sumOfProb(lowPrev)
	assert.String(sop != 0, "no following weights to "+lowPrev)
	extracted := uint(rand.Uint64()) % sop
	for word, prob := range m[lowPrev] {
		if extracted < prob {
			return word
		}
		extracted -= prob
	}
	assert.String(false, "extracted number out of range")
	return ""
}

func (m Markov) Generate() string {
	var res, prev string
	for {
		word := m.genWord(prev)
		if word == "" {
			break
		}
		res += word + " "
		prev = word
	}
	// return strings.Trim(res, " ")
	return res
}
