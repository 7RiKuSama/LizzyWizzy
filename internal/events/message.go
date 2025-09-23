package events

import (
	"log"
	"strings" // Add this import

	"github.com/bwmarrin/discordgo"
)

func (h *Handlers) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if s.State.User.ID == m.Author.ID {
        return
    }

    // Trim whitespace from the message content
    content := strings.TrimSpace(m.Content)

    // Check for the prefix before doing anything else
    if len(content) > 0 && content[0] != '!' {
        if err := h.Services.IncrementExp(h.Context, m); err != nil {
            log.Println(err)
        }
    }

    if len(content) > 0 {
        switch content {
        case "!play":
            Player.JoinVoiceChannel(s, m)
        case "!leave":
            Player.LeaveVoiceChannel(s, m)
        case "!next":
            Player.NextTrack(s, m, true)
        case "!previous":
            Player.PreviousTrack(s, m)
        case "!rank":
            h.Services.LoadRank(h.Context, s, m)
        case "!nowplaying":
            Player.MusicInfo(s, m)
        }
    }
}
