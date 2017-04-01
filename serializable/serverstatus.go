package serializable

// Holds the Server Status.
type ServerStatus struct {
	Version Version `json:"version,string"`
	Players Players `json:"players,string"`
	Description string `json:"description"`
	Favicon string `json:"favicon"`
}

// Holds the status Versions.
type Version struct {
	Name string `json:"name"`
	Protocol int `json:"protocol"`
}

// Holds the status Players.
type Players struct {
	Max int `json:"max"`
	Online int `json:"online"`
	Sample []Player `json:"sample"`
}

// Holds the status Player.
type Player struct {
	Name string `json:"name"`
	Id string `json:"id"`
}