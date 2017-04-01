package networking

import (
	"GoMiniServer/chunk"
	"GoMiniServer/networking/buffer"
	"github.com/satori/go.uuid"
	"math/rand"
)

//TODO: Fully restructure this to have read and write methods.
type PacketHandler interface {
	ID() byte
	Handle(miniServer *MiniServer, in *InboundPacket, out *buffer.MiniBuffer)
}

// A map of all the Packets in order to be retrievable in the future.
var PacketHandlers = map[uint16]map[byte]PacketHandler{

	316: {
		0x00: HandShakeHandle{},
		0x01: PingHandle{},
		0x09: LoginSuccessHandle{},
	},
}

// Handshake Handler.
type HandShakeHandle struct{}

func (packet HandShakeHandle) ID() byte { return 0x00 }

func (packet HandShakeHandle) Handle(miniServer *MiniServer, in *InboundPacket, out *buffer.MiniBuffer) {

	miniBuf := in.Buffer

	if in.PacketLength >= 16 {

		var state int
		println("Protocol:", miniBuf.ReadVarInt(), "IP:", miniBuf.ReadString(), "Port:", miniBuf.ReadUnsignedShort(), "State:", func() int { state = miniBuf.ReadVarInt(); return state }())

		out.WriteVarInt(0x00)
		if state == 1 {
			out.WriteString(miniServer.SerializedStatus)
		} else {
			out.WriteString("{\"text\":\"§c§lSworreh, §r§awe will be back soon, with some fanceh developments!\"}")
		}

		//Login Start --> Login Success Skips Encryption
	} else {

		in.Player.UUID = uuid.NewV4()
		in.Player.Username = in.Buffer.ReadString()

		out.WriteVarInt(0x02)
		out.WriteString(in.Player.UUID.String())
		out.WriteString(in.Player.Username)

		out.WriteTo(in.Player.Connection)
		out.ClearAll()

		println("Username:", in.Player.Username, "UUID:", in.Player.UUID.String())

		// Next: Send http://wiki.vg/Protocol#Join_Game

		out.WriteVarInt(0x23)

		// Sets a Random EID (Entity ID)
		out.WriteInt(rand.Int())

		// Set Gamemode to Survival (0)
		out.WriteBytes(0)

		// Set Dimension to OverWorld (0)
		out.WriteInt(0)

		// Set Difficulty to Peaceful (0) and Sets Max players to 100
		out.WriteBytes(0, 100)

		// Sets the Level type to Default.
		out.WriteString("default")

		// Sets Reduced Debug Info to false by sending a bit of value 0 (Might not work.)
		out.WriteBytes(0)

		/*println(in.Buffer.ReadString())
		out.WriteVarInt(0x01)
		out.WriteString("MiniServer")
		out.Write(miniServer.PubKey)
		out.Write(miniServer.SecretKey)
		in.Player.HasEncryption = true*/
	}
}

type PluginMessageHandle struct{}

func (packet PluginMessageHandle) ID() byte { return 0x01 }

// Handles PluginMessage Response
func (packet PluginMessageHandle) Handle(miniServer *MiniServer, in *InboundPacket, out *buffer.MiniBuffer) {
	//out.WriteVarInt(0x0C)
}

// Handles normal response to ping.
type PingHandle struct{}

func (packet PingHandle) ID() byte { return 0x01 }

// Handles Ping response.
func (packet PingHandle) Handle(miniServer *MiniServer, in *InboundPacket, out *buffer.MiniBuffer) {
	in.Buffer.WriteTo(in.Player.Connection)
	println("Ping")
}

// TODO: Decrypt packets, by making a Player struct which stores the key.
type LoginSuccessHandle struct{}

func (packet LoginSuccessHandle) ID() byte { return 0x04 }

func (packet LoginSuccessHandle) Handle(miniServer *MiniServer, in *InboundPacket, out *buffer.MiniBuffer) {

	// Send Spawn position
	// out.WriteVarInt(0x43)

	// // Generates a position with 0, 256, 0
	// out.WriteInt(((0 & 0x3FFFFFF) << 38) | ((256 & 0xFFF) << 26) | (0 & 0x3FFFFFF))

	// // Sends to the Connection
	// out.WriteTo(in.Player.Connection)

	// out.ClearAll()

	// Send Chunks
	out.WriteVarInt(0x20)

	// Send Chunk X
	out.WriteInt(0)

	// Send Chunk Z
	out.WriteInt(0)

	// Ground-Up Continuous (true)
	out.WriteBytes(1)

	out.WriteVarInt(0)

	section := chunk.Chunk{X:0, Z:0}

	section.Fill()

	compiled := ""
	for _, value := range section.Sections { compiled += value.Blocks }

	println(compiled)
	out.WriteString(compiled)


	// file, err := os.Open("bigtest.nbt")
	// if err != nil { log.Fatal(err) }

	// reader, err := gzip.NewReader(file)
	// if err != nil { log.Fatal(err) }

	// nbtStuff, _ := nbt.Decode(reader)
	// nbtStuff.PrettyPrint()

	// out.Write

	// out.WriteVarInt(0)

	//out.WriteVarInt()

	println("Sent!")
}
