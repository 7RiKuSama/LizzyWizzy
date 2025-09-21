package commands

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/7RikuSama/liz.git/internal/response"
	"github.com/7RikuSama/liz.git/internal/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/dhowden/tag"
)

type Player struct {
	IsPlaying    bool
	IsSlash      bool
	TrackNumber  int
	TotalTracks  int
	Files        []string
	Raw          *os.File
}

func NewPlayer(s *discordgo.Session, m *discordgo.MessageCreate, i *discordgo.InteractionCreate) *Player {
	return &Player{
		TrackNumber:  0,
		TotalTracks:  0,
		Files:        []string{},
		Raw:          nil,
	}
}

func (p *Player) SetResponse(session *discordgo.Session, event any, response *discordgo.MessageSend) {
	switch v := event.(type) {
	case *discordgo.MessageCreate:
		session.ChannelMessageSendComplex(v.ChannelID, response)
	case *discordgo.InteractionCreate:
		session.InteractionRespond(v.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: response.Embeds,
			},
		})
	}
}

type DeezerTrack struct {
	Data []struct {
		Album struct {
			Cover string `json:"cover"`
		}
	} `json:"data"`
}

func (p *Player) AddMusicFiles() error {
	files, err := os.ReadDir("assets/music")
	if err != nil {
		return err
	}
	p.TotalTracks = len(files)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		p.Files = append(p.Files, file.Name())
	}
	return nil
}

func (p *Player) GetFile() error {
	file, err := os.Open(fmt.Sprintf("assets/music/%s", p.Files[p.TrackNumber]))

	if err != nil {
		return err
	}

	p.Raw = file
	return nil
}

func (p *Player) MusicEmbed(title string) *discordgo.MessageSend {
	// --- Metadata fields ---
	var fields []*discordgo.MessageEmbedField
	var trackTitle, artist, album string

	// Make sure Raw is not nil
	if p.Raw != nil {
		if m, err := tag.ReadFrom(p.Raw); err == nil {
			trackTitle = m.Title()
			artist = m.Artist()
			album = m.Album()

			fields = []*discordgo.MessageEmbedField{
				{
					Name:   "🎵 Title",
					Value:  trackTitle,
					Inline: true,
				},
				{
					Name:   "👤 Artist",
					Value:  artist,
					Inline: false,
				},
				{
					Name:   "💿 Album",
					Value:  album,
					Inline: true,
				},
			}
		}
	}

	// If title/artist is missing, fall back
	if trackTitle == "" && p.TrackNumber < len(p.Files) {
		trackTitle = p.Files[p.TrackNumber]
	}
	if artist == "" {
		artist = "Unknown Artist"
	}
	if album == "" {
		album = "Unknown Album"
	}

	cover := "https://media2.giphy.com/media/v1.Y2lkPTc5MGI3NjExeWVja3o1M2F5dmR4cmg4Nms4eTE1OWFtMHFva2FtdnF1OXc2bG51YSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/zLjVqZQx3JNl8vOwk6/giphy.gif"

	dataUrl := "https://api.deezer.com/search"
	query := url.Values{}
	query.Set("q", fmt.Sprintf("%s %s", artist, trackTitle))
	query.Set("limit", "1")

	if resp, err := http.Get(dataUrl + "?" + query.Encode()); err == nil && resp != nil {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err == nil {
			var result DeezerTrack
			if json.Unmarshal(body, &result) == nil && len(result.Data) > 0 {
				cover = result.Data[0].Album.Cover
			}
		}
	}

	// --- Spotify Search Link ---
	spotifyURL := fmt.Sprintf("https://open.spotify.com/search/%s",
		url.QueryEscape(trackTitle+" "+artist))

	// --- Build Embed ---
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					IconURL: "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcdnl.iconscout.com%2Flottie%2Fpremium%2Fthumb%2Faudio-wave-5692234-4798534.gif&f=1&nofb=1&ipt=1d0784fd68e98da166a187a8c3d3318410ea88df75c9018c3102b9e7670ec423",
					Name:    "Now Playing",
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL:    cover,
					Height: 500,
					Width:  500,
				},
				Title:       title,
				Description: fmt.Sprintf(response.MessageNowPlaying, trackTitle),
				Color:       0xe972bd,
				Fields:      fields,
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Search on Spotify",
						Style: discordgo.LinkButton,
						URL:   spotifyURL,
					},
				},
			},
		},
	}
}

func (p *Player) NextTrack(session *discordgo.Session, event any, userIssued bool) {
	if !p.IsPlaying {
		p.SetResponse(session, event, utils.SetWarningMessage("Next", response.MessageNoVoiceChannel))
		return
	}
	if p.TotalTracks > 0 {
		p.TrackNumber = (p.TrackNumber + 1) % p.TotalTracks
	}

	if err := p.GetFile(); err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Next", response.ErrNextTrack))
		return
	}

	if userIssued {
		p.SetResponse(session, event, p.MusicEmbed("Next"))
	}
}

func (p *Player) PreviousTrack(session *discordgo.Session, event any) {
	if !p.IsPlaying {
		p.SetResponse(session, event, utils.SetWarningMessage("Previous", response.MessageNoVoiceChannel))
		return
	}
	if p.TotalTracks > 0 {
		if p.TrackNumber > 0 {
			p.TrackNumber--
		} else {
			p.TrackNumber = p.TotalTracks - 1
		}
	}
	if err := p.GetFile(); err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Previous", response.ErrPreviousTrack))
		return
	}
	p.SetResponse(session, event, p.MusicEmbed("Previous"))
}

// function that returns the channelID of that the user is joining
func (p *Player) getVoiceChannel(session *discordgo.Session, event any) (string, error) {
	var author, guild string
	switch v := event.(type) {
	case *discordgo.InteractionCreate:
		author = v.Member.User.ID
		guild = v.GuildID
	case *discordgo.MessageCreate:
		author = v.Author.ID
		guild = v.GuildID
	}

	vs, err := session.State.VoiceState(guild, author)

	if err != nil {
		return "", err
	}

	return vs.ChannelID, nil
}

// Command to make the bot play music in the voice channel that the user is in
func (p *Player) JoinVoiceChannel(session *discordgo.Session, event any) {

	var guild string

	switch v := event.(type) {
	case *discordgo.MessageCreate:
		guild = v.GuildID
	case *discordgo.InteractionCreate:
		guild = v.GuildID
	}

	vc := session.VoiceConnections[guild]

	if vc != nil && vc.ChannelID != "" {

		return
	}

	if err := p.AddMusicFiles(); err != nil {
		panic("Coudln't read the music directory")
	}

	channelID, err := p.getVoiceChannel(session, event)

	if err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Play", response.ErrUnableToJoinVC))
		return
	}

	vc, err = session.ChannelVoiceJoin(guild, channelID, false, false)

	if err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Play", response.ErrUnableToJoinVC))
		return
	}

	defer vc.Disconnect()
	if err := p.GetFile(); err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Error", response.ErrPlayError))
		return
	}

	p.SetResponse(session, event, p.MusicEmbed("Playing"))

	index := p.TrackNumber
	p.IsPlaying = true
	for {

		index = index % len(p.Files)

		ffmpeg := exec.Command("ffmpeg", "-i", "assets/music/"+p.Files[index], "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
		dca := exec.Command("dca", "-as", "960")

		pipe, err := ffmpeg.StdoutPipe()
		if err != nil {
			p.SetResponse(session, event, utils.SetErrorMessage("Error", response.ErrPlayError))
			continue
		}
		dca.Stdin = pipe

		output, err := dca.StdoutPipe()
		if err != nil {
			p.SetResponse(session, event, utils.SetErrorMessage("Error", response.ErrPlayError))
			log.Println(err)
			continue
		}

		if err := ffmpeg.Start(); err != nil {
			p.SetResponse(session, event, utils.SetErrorMessage("Error", response.ErrPlayError))
			log.Println(err)
			continue
		}
		if err := dca.Start(); err != nil {
			p.SetResponse(session, event, utils.SetErrorMessage("Error", response.ErrPlayError))
			log.Println(err)
			continue
		}

		vc.Speaking(true)

		// Read opus frames safely
		curr := index
		for {
			if curr != p.TrackNumber {
				index = p.TrackNumber
				break
			}
			var opusLen int16
			if err := binary.Read(output, binary.LittleEndian, &opusLen); err != nil {
				break // EOF or pipe closed
			}
			opus := make([]byte, opusLen)
			if _, err := io.ReadFull(output, opus); err != nil {
				break
			}
			vc.OpusSend <- opus
		}

		if curr != p.TrackNumber {
			continue
		}
		p.NextTrack(session, event, false)
		ffmpeg.Wait()
		dca.Wait()
	}
}

func (p *Player) LeaveVoiceChannel(session *discordgo.Session, event any) {
    var guild string
    switch v := event.(type) {
    case *discordgo.MessageCreate:
        guild = v.GuildID
    case *discordgo.InteractionCreate:
        guild = v.GuildID
    }

    // Check if a voice connection exists for the guild.
    // If it doesn't, inform the user and return.
    if session.VoiceConnections[guild] == nil {
        p.SetResponse(session, event, utils.SendMessage("Leaving", response.MessageNoVoiceChannel)) // Use an appropriate message
        return
    }

    session.VoiceConnections[guild].Disconnect()
    p.IsPlaying = false
    p.SetResponse(session, event, utils.SendMessage("Leaving", response.MessageLeaveChannel))
}

func (p *Player) MusicInfo(session *discordgo.Session, event any) {
	if !p.IsPlaying {
		p.SetResponse(session, event, utils.SetWarningMessage("Not Connected", response.MessageNoVoiceChannel))
		return
	}

	if err := p.GetFile(); err != nil {
		p.SetResponse(session, event, utils.SetErrorMessage("Error", "Could not load track file"))
		return
	}
	p.SetResponse(session, event, p.MusicEmbed("Song"))
}

func (p *Player) SetFavorite(slash bool) {

}
