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
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/sgorblex/MarkovGenBot/backend"
)

func main() {
	markovs := make(Tables)

	// save training to disk and exit gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Fprintf(os.Stderr, "\nexiting gracefully...\n")
		GracefulExit(markovs)
		os.Exit(0)
	}()

	updates := GetUpdates()
	for update := range updates {
		err := ProcessUpdate(markovs, update)
		if err != nil {
			log.Println(err)
		}
	}
}
