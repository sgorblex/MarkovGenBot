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
	var prev, word string
	for _, word = range strings.Fields(text) {
		if m[prev] == nil {
			m[prev] = make(Following)
		}
		m[prev][word]++
		prev = word
	}
	if m[word] == nil {
		m[word] = make(Following)
	}
	m[word][""]++
}

func (m Markov) sumOfProb(word string) uint {
	var i uint
	for _, prob := range m[word] {
		i += prob
	}
	return i
}

func (m Markov) genWord(prev string) string {
	sop := m.sumOfProb(prev)
	assert.String(sop != 0, "no following weights to "+prev)
	extracted := uint(rand.Uint64()) % sop
	for suff, prob := range m[prev] {
		if extracted < prob {
			return suff
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
