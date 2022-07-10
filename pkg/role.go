package pkg

import (
	"fmt"
	"github.com/GuessWhoSamFoo/gang-gang-bot/pkg/util"
	"golang.org/x/exp/slices"
)

type FieldType string

const (
	AcceptedIcon  = "✅"
	DeclinedIcon  = "❌"
	TentativeIcon = "❔"

	AcceptedField  FieldType = "Accepted"
	DeclinedField  FieldType = "Declined"
	TentativeField FieldType = "Tentative"
	WaitlistField  FieldType = "Waitlist"
)

type Role struct {
	Icon      string
	FieldName FieldType
	Users     []string
	Count     int
	Limit     int
}

// NewRole creates an instance of a role
func NewRole(icon string, fieldName FieldType) *Role {
	return &Role{
		Icon:      icon,
		FieldName: fieldName,
		Users:     []string{},
	}
}

type RoleGroup struct {
	Roles    []*Role
	Waitlist map[FieldType]*Role
}

// NewDefaultRoleGroup returns the default Accepted, Declined, and Tentative fields
func NewDefaultRoleGroup() *RoleGroup {
	rg := &RoleGroup{
		Roles:    []*Role{},
		Waitlist: map[FieldType]*Role{},
	}
	rg.AddRole(
		NewRole(AcceptedIcon, AcceptedField),
		NewRole(DeclinedIcon, DeclinedField),
		NewRole(TentativeIcon, TentativeField),
	)
	rg.AddWaitlistForRole("", AcceptedField)
	return rg
}

func (rg *RoleGroup) AddRole(roles ...*Role) {
	rg.Roles = append(rg.Roles, roles...)
}

func (rg *RoleGroup) AddWaitlistForRole(icon string, field FieldType) {
	rg.Waitlist[field] = &Role{
		Icon:      icon,
		FieldName: WaitlistField,
		Users:     []string{},
	}
}

func (rg *RoleGroup) ToggleRole(fieldName FieldType, user string) error {
	for _, r := range rg.Roles {
		hasUser := slices.Contains(r.Users, user)
		isFull := r.Limit > 0 && r.Count == r.Limit
		wl, hasWaitlist := rg.Waitlist[r.FieldName]
		switch {
		case hasUser && hasWaitlist:
			r.Count--
			r.Users = util.RemoveUser(r.Users, user)
			if wl.Count > 0 {
				wl.Count--
				name := wl.Users[0]
				wl.Users = wl.Users[1:]
				r.Count++
				r.Users = append(r.Users, name)
			}
		case hasUser && !hasWaitlist:
			r.Count--
			r.Users = util.RemoveUser(r.Users, user)
		case !hasUser && hasWaitlist:
			if slices.Contains(wl.Users, user) {
				wl.Count--
				wl.Users = util.RemoveUser(wl.Users, user)
				continue
			}
			if isFull && r.FieldName == fieldName {
				wl.Count++
				wl.Users = append(wl.Users, user)
				continue
			}
			if r.FieldName == fieldName {
				r.Count++
				r.Users = append(r.Users, user)
			}
		case !hasUser && !hasWaitlist:
			if isFull {
				return fmt.Errorf("event full; cannot add to waitlist")
			}
			if r.FieldName == fieldName {
				r.Count++
				r.Users = append(r.Users, user)
			}
		}
	}
	return nil
}

// RemoveFromAllLists removes a username from all role groups including waitlists
func (rg *RoleGroup) RemoveFromAllLists(name string) error {
	for _, r := range rg.Roles {
		r.Users = util.RemoveUser(r.Users, name)
		r.Count = len(r.Users)
		wl, ok := rg.Waitlist[r.FieldName]
		if ok {
			wl.Users = util.RemoveUser(wl.Users, name)
			wl.Count = len(wl.Users)
		}
	}
	return nil
}

// PeekWaitlist returns the first user on the waitlist, if any
func (rg *RoleGroup) PeekWaitlist(field FieldType) string {
	wl, ok := rg.Waitlist[field]
	if ok && wl.Count > 0 {
		return wl.Users[0]
	}
	return ""
}

// HasUser checks if a given role has a user
func (rg *RoleGroup) HasUser(name string, field FieldType) bool {
	for _, r := range rg.Roles {
		if r.FieldName == field && slices.Contains(r.Users, name) {
			return true
		}
	}
	return false
}
