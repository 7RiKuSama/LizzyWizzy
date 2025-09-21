package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *Handlers) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	Player.Message = m

	if s.State.User.ID == m.Author.ID {
		return
	}

	err := h.Services.CreateMembership(h.Context, m)
	if err != nil {
		log.Println(err)
	}

	if m.Content[0] != '!' {
		if err := h.Services.IncrementExp(h.Context, m); err != nil {
			log.Println(err)
		}
	}

	if err := h.Services.LevelTracker(h.Context, s, m); err != nil {
		log.Println(err)
	}

	switch m.Content {
	case "!play":
		Player.JoinVoiceChannel(false)
	case "!leave":
		Player.LeaveVoiceChannel(false)
	case "!next":
		Player.NextTrack(false, true)
	case "!previous":
		Player.PreviousTrack(false)
	case "!rank":
		h.Services.LoadRank(h.Context, s, m)
	case "!nowplaying":
		Player.MusicInfo(false)
	}
}
