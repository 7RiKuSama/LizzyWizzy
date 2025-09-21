package response

const (
	// Client Errors
	ErrPlayError       = "Signal corrupted — my voice fractured in the stream. The track can’t go on…"
	ErrUnableToJoinVC  = "Connection refused. My signal can’t reach your channel — the feed’s broken, darling. 🔧"
	ErrNoTokenProvided = "Please Provide a valid Bot Token"
	ErrNoChannelID     = "Invaid Channel ID"
	ErrNextTrack       = "Mmh… I tried to spin the next track, but it wouldn’t play. Maybe it’s broken, maybe it’s shy. 🎶💔"
	ErrPreviousTrack   = "I tried to take you back, but the last track is lost in silence. Some memories can’t be replayed, darling. 🎶"
)
