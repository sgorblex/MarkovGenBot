package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/sgorblex/MarkovGenBot/markov"
)

func MarkovsFromFile(filePath string) (map[ChatID]markov.Markov, error) {
	if _, err := os.Stat(filePath); err != nil {
		return make(map[ChatID]markov.Markov), nil
	}
	file, err := os.Open(filePath)
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

func GracefulExit(m map[ChatID]markov.Markov, filePath string) {
	fmt.Fprintf(os.Stderr, "\nexiting gracefully...\n")
	jsonData, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
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
