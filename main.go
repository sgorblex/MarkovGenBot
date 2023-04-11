package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/sgorblex/MarkovGenBot/backend"
)

const (
	persistTimer time.Duration = 600
)

func main() {
	markovs := make(Tables)

	// save training to disk and exit gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exiting gracefully")
		Persist(markovs)
		os.Exit(0)
	}()

	// periodically write tables to persistence
	go func() {
		ticker := time.NewTicker(persistTimer * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			log.Println("Executing persistence routine")
			Persist(markovs)
		}
	}()

	updates := GetUpdates()
	for update := range updates {
		err := ProcessUpdate(markovs, update)
		if err != nil {
			log.Println(err)
		}
	}
}
