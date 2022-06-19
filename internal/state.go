package internal

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

type StateManager struct {
	ActiveMap
	CommandHandlers   map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewStateManager() *StateManager {
	sm := &StateManager{}
	sm.ActiveMap = ActiveMap{
		userMap: make(map[string]struct{}),
	}
	sm.CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"event": sm.CreateEventHandler,
		//"my_events": ListEventHandler,
		//"edit":      EditEventHandler,
	}
	sm.ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"accept":        sm.AcceptHandler,
		"decline":       sm.DeclineHandler,
		"tentative":     sm.TentativeHandler,
		"edit":          sm.EditHandler,
		"delete":        sm.DeleteHandler,
		"confirmDelete": sm.ConfirmDeleteHandler,
	}
	return sm
}

type ActiveMap struct {
	mu      sync.Mutex
	userMap map[string]struct{}
}

func (a *ActiveMap) AddUser(user string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.userMap[user] = struct{}{}
}

func (a *ActiveMap) RemoveUser(user string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.userMap, user)
}

func (a *ActiveMap) HasUser(user string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.userMap[user]
	return ok
}
