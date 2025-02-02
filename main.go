package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/NotNotQuinn/go-irc/channels"
	"github.com/NotNotQuinn/go-irc/client"
	"github.com/NotNotQuinn/go-irc/cmd"
	"github.com/NotNotQuinn/go-irc/config"
	"github.com/NotNotQuinn/go-irc/core/incoming"
	"github.com/NotNotQuinn/go-irc/core/sender"
	"github.com/NotNotQuinn/go-irc/handlers"
)

func main() {
	defer recoverFromDisconnect()
	go handleErrors()
	go incoming.HandleAll()

	fmt.Print("Starting")
	err := config.Init()
	if err != nil {
		panic(err)
	}

	// Dots to show progress, even though they mostly go all at once
	// its a good measure of startup speed changing over time.
	fmt.Print(".")
	cmd.LoadAll()

	fmt.Print(".")
	cc, err := client.GetCollection()
	if err != nil {
		panic(err)
	}

	fmt.Print(".")
	handlers.Handle(cc)

	fmt.Print(".")
	err = cc.JoinAll()
	if err != nil {
		panic(err)
	}
	go sender.HandleAllSends(cc)

	fmt.Print(".")
	err = cc.Connect()
	if err != nil {
		panic(err)
	}
}

// Handles all errors
func handleErrors() {
	for {
		// although it doesnt seem like much, it allows for good error logging later on.
		// Errors should only be passed to this stream if there is no other place, and
		// a panic is not sutible
		err := <-channels.Errors
		fmt.Printf("Error: %+v\n", err)
	}
}

// Increases as restart attempts increace in count
var restartMult = 1

// Attempts to recover from a disconnect, re-panics other errors
func recoverFromDisconnect() {
	if err := recover(); err != nil {
		s := fmt.Sprint(err)
		if strings.Contains(s, "no such host") && strings.Contains(s, "irc.chat.twitch.tv") {
			if !(restartMult >= 32) {
				// will never exceed 32
				// max amount of time waited is 8 mins (15 * 2^5 seconds)
				restartMult *= 2
			}
			sleepTime := time.Second * 15 * time.Duration(restartMult)
			fmt.Println("\nConnection interupted, attempting restart in", sleepTime)
			time.Sleep(sleepTime)
			main()
		}
		panic(err)
	}
}
