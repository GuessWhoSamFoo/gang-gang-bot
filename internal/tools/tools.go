package main

import (
	"github.com/GuessWhoSamFoo/fsm"
	"github.com/GuessWhoSamFoo/gang-gang-bot/internal/commands"
	"log"
	"os"
)

func main() {
	f := fsm.NewFSM("", commands.CreateEvents(), nil)
	createGraph := fsm.Visualize(f)
	if err := os.WriteFile("fsm_create.dot", []byte(createGraph), 0600); err != nil {
		log.Fatal(err)
	}

	editGraph := fsm.Visualize(fsm.NewFSM("", commands.EditEvents(), nil))
	if err := os.WriteFile("fsm_edit.dot", []byte(editGraph), 0600); err != nil {
		log.Fatal(err)
	}
}
