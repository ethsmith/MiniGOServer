package networking

import (
	"net"
	"strconv"
	"GoMiniServer/networking/buffer"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"GoMiniServer/player"
	"io"
	"github.com/jinzhu/gorm"
)

// Don't recommend directly initializing.
type MiniServer struct {

	gorm.Model

	Ip               string
	Port             int
	Listener         *net.TCPListener `json:"-"`
	SerializedStatus string
	PrivKey          *rsa.PrivateKey
	PubKey           []byte
	SecretKey        []byte
	//EventBus         (*EventBus)
}

// A Incoming packet.
type InboundPacket struct {
	Buffer           *buffer.MiniBuffer
	Player           *player.Player
	PacketLength     int
}

// Creates a new MiniServer instance with requested perimeters.
func NewMiniServer(ip string, port int, serializedStatus string) *MiniServer {

	secretKeyBuffer := buffer.NewMiniBuffer(make([]byte, 0))
	io.CopyN(secretKeyBuffer, rand.Reader, 4)

	privKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	pubKey, _ := x509.MarshalPKIXPublicKey(&privKey.PublicKey)

	return &MiniServer{Ip: ip, Port: port, PrivKey: privKey, PubKey: pubKey, SecretKey: secretKeyBuffer.Bytes, SerializedStatus: serializedStatus}
}

func (miniServer *MiniServer) Enabled() bool {
	return miniServer.Listener != nil
}


// Handles the incoming connection.
func handlePacket(miniServer *MiniServer, packetLength int, player *player.Player, in *buffer.MiniBuffer, out *buffer.MiniBuffer) {

	if packetLength == 1 { return }

	id := in.ReadVarInt()

	println("Length:", packetLength, "ID:", "0x" + strconv.FormatInt(int64(id), 16))

	packet := PacketHandlers[316][byte(id)]

	if packet == nil { println("Packet not found", id); return }

	packet.Handle(miniServer, &InboundPacket { Buffer: in, Player: player, PacketLength: packetLength }, out)
}

// Starts the server, returns false if the server is already enabled or has failed to start.
func (miniServer *MiniServer) Start() bool {

	if miniServer.Enabled() { return false }

	go miniServer.startListener()

	return true
}

// Stops the server.
func (miniServer *MiniServer) Stop() {

	if !miniServer.Enabled() { return }

	miniServer.Listener.Close()

	miniServer.Listener = nil

}

// Starts the listener, which listens for TCP connections and directs them accordingly.
func (miniServer *MiniServer) startListener() {

	tcpIp, _ := net.ResolveTCPAddr("tcp", miniServer.Ip + ":" + strconv.Itoa(miniServer.Port))

	miniServer.Listener, _ = net.ListenTCP("tcp", tcpIp)

	for miniServer.Enabled() {

		connection, _ := miniServer.Listener.AcceptTCP()

		connection.SetNoDelay(true)
		connection.SetKeepAlive(true)

		p := player.Player{Connection: connection, HasEncryption: false}

		go func() {

			out := buffer.NewMiniBuffer(make([]byte, 0))

			println("Connection from:", connection.RemoteAddr().String())

			for {

				packetInfo := make([]byte, 1)
				length, err := connection.Read(packetInfo)

				if err != nil || length == 0 { println("Finished reading."); println(); break }

				packetLength := int(packetInfo[0])

				readBytes := make([]byte, packetLength)

				connection.Read(readBytes)

				handlePacket(miniServer, packetLength, &p, buffer.NewMiniBuffer(readBytes), out)

				if len(out.Bytes) != 0 { out.WriteTo(connection); out.ClearAll(); println("Sent something") }

			}

		}()

	}

}