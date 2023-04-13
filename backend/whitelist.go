package backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	whitelistsPath string = "whitelist.json"
)

type empty struct{}

var whitelist map[ChatID]empty
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
	var whitelistList []ChatID
	err = json.Unmarshal(raw, &whitelistList)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Parsed whitelist: %v.\n", whitelistList)
	whitelist = make(map[ChatID]empty, len(whitelistList))
	for _, w := range whitelistList {
		whitelist[w] = empty{}
	}
	useWhitelist = true
}

func isWhitelisted(cID ChatID) bool {
	_, isIn := whitelist[cID]
	return isIn
}
