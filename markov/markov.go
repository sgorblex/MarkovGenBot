package markov

// TODO
// see comments with TODO

import (
	"log"
	"math/rand"
	"strings"
)

type Following map[string]uint
type Markov map[string]Following

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

func (m Markov) sumOfProb(word string) uint {
	var i uint
	for _, prob := range m[word] {
		i += prob
	}
	return i
}

func (m Markov) genWord(word string) string {
	// TODO: fix this
	sop := m.sumOfProb(word)
	if sop == 0 {
		return "Failed. Fix this bug, bro."
	}
	extracted := uint(rand.Uint64()) % sop
	// TODO: check algorithm and condition
	for suff, prob := range m[word] {
		extracted -= prob
		if extracted < 0 {
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
	// TODO: make this selection vaguely coherent with the genWord one
	r = rand.Intn(len(m))
	for prev = range m {
		if r >= j {
			break
		}
		j++
	}
	for i := 0; i < nWords; i++ {
		word := m.genWord(prev)
		res += word + " "
		prev = word
	}
	return res
}
