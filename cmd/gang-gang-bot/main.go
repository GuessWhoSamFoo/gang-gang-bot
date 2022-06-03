package main

import (
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := internal.NewConfig()
	if err != nil {
		log.Fatalln("Failed to get config file")
	}
	bot, err := services.NewBot(config)
	if err != nil {
		log.Fatalln("Failed initialize bot")
	}
	defer func() {
		if err := bot.Close(); err != nil {
			log.Fatalf("cannot cleanup session: %v", err)
		}
	}()

	if err := bot.Start(); err != nil {
		log.Fatalf("cannot start bot: %v", err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGKILL, os.Interrupt, os.Kill)
	<-stop
	log.Println("Goodbye")
}
