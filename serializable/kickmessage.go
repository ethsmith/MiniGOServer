package serializable

// Helps people use the right format when kicking a player.
type KickMessage struct {
	Text string `json:"text"`
}

