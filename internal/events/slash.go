package events

import (
	"log"

	"github.com/7RikuSama/liz.git/internal/commands"
	"github.com/bwmarrin/discordgo"
)

type Slash struct {
	Commands        map[string]*discordgo.ApplicationCommand
	CommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewSlash() *Slash {
	sl := &Slash{
		Commands:        make(map[string]*discordgo.ApplicationCommand),
		CommandHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}
	return sl
}

// this function initialize the commands, and their respective handlers in a slash
func (h *Handlers) initCommands() {
	h.Slash.Commands["rank"] = &discordgo.ApplicationCommand{
		Name:        "rank",
		Description: "display user's rank",
	}

	h.Slash.Commands["play"] = &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "joins the voice channel and start playing music",
	}

	h.Slash.Commands["leave"] = &discordgo.ApplicationCommand{
		Name:        "leave",
		Description: "leaves the voice channel",
	}

	h.Slash.Commands["nowplaying"] = &discordgo.ApplicationCommand{
		Name:        "nowplaying",
		Description: "displays info about the current song",
	}

	h.Slash.CommandHandlers["rank"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := commands.LoadRank(h.Context, h.Session, h.Services, i)
		if err != nil {
			log.Println("err: ", err)
		}
	}

	h.Slash.CommandHandlers["play"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		Player.JoinVoiceChannel(h.Session, i)
	}

	h.Slash.CommandHandlers["leave"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		Player.LeaveVoiceChannel(h.Session, i)
	}

	h.Slash.CommandHandlers["nowplaying"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		Player.MusicInfo(h.Session, i)
	}
}

// AddCommands register slash commands
func (h *Handlers) AddCommands(session *discordgo.Session, guildID string) {
	for key, v := range h.Slash.Commands {
		cmd, err := session.ApplicationCommandCreate(session.State.Application.ID, guildID, v)
		if err != nil {
			panic(err)
		}
		h.Slash.Commands[key] = cmd
	}
}

// DeleteCommands deletes all the registered slash commands
func (h *Handlers) DeleteCommands(session *discordgo.Session, guildID string) {
	for _, cmd := range h.Slash.Commands {
		err := session.ApplicationCommandDelete(session.State.Application.ID, guildID, cmd.ID)
		if err != nil {
			panic(err)
		}
	}
}

// OnInteractionCreate handles income slash commands interactions
func (h *Handlers) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if handler, ok := h.Slash.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
				},
			})
			handler(s, i)
		} else {
			// This block handles unknown commands and prevents a timeout.
			log.Println("Unknown command received:", i.ApplicationCommandData().Name)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sorry, I don't know that command!",
				},
			})
		}
	}
}
