package player

import (
	"net"
	"github.com/satori/go.uuid"
)

// Official Player Struct.
type Player struct {
	Connection       *net.TCPConn
	Username         string
	UUID             uuid.UUID
	SecretKey        string
	HasEncryption    bool
}