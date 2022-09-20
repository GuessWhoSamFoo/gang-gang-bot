package pkg

import "github.com/bwmarrin/discordgo"

var (
	EventPermission int64 = discordgo.PermissionManageEvents
	DMPermission          = true
	// Privileged Gateway Intents for server members must be enabled on https://discord.com/developers/applications
)
