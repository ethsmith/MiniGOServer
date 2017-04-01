package buffer

import (
	"net"
	"io"
)


// TODO: Add Read/Write VarLong functions

// A Minecraft regulated Buffer.
// Don't recommend directly initializing.
type MiniBuffer struct {
	io.Writer

	Bytes []byte
	ReaderIndex int
	WriterIndex int
}


// Returns a new MiniBuffer with a provided Array.
func NewMiniBuffer(bytes []byte) *MiniBuffer {
	return &MiniBuffer{Bytes: bytes}
}

// Returns a new MiniBuffer with a provided Array, readerIndex and writerIndex.
func NewMiniBufferWithIndex(bytes []byte, readerIndex int) *MiniBuffer {
	return &MiniBuffer{Bytes: bytes, ReaderIndex: readerIndex}
}

func (miniBuffer *MiniBuffer) WriteTo(connection net.Conn) {

	length := len(miniBuffer.Bytes)
	wrapper := NewMiniBuffer(make([]byte, length))

	wrapper.WriteVarInt(length)
	wrapper.WriteBytes(miniBuffer.Bytes...)

	connection.Write(wrapper.Bytes)

}

// Clears all bytes in the backend array.
func (miniBuffer *MiniBuffer) ClearAll() {
	miniBuffer.ClearBeyond(0)
}

// Clears beyond a certain index.
func (miniBuffer *MiniBuffer) ClearBeyond(index int) {
	miniBuffer.Bytes = miniBuffer.Bytes[:index]
	if index < miniBuffer.ReaderIndex { miniBuffer.ReaderIndex = index }
	if index < miniBuffer.WriterIndex { miniBuffer.WriterIndex = index }
}


// Writes Bytes to the end of the backend buffer.
func (miniBuffer *MiniBuffer) WriteBytes(b... byte) {
	miniBuffer.Bytes = append(miniBuffer.Bytes, b...)
	miniBuffer.WriterIndex += len(b)
}

// Writes Short in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteShort(short int16) {
	miniBuffer.WriteBytes(byte(short >> 8), byte(short & 0xff))
}

// Writes Unsigned Short in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteUnsignedShort(short uint16) {
	miniBuffer.WriteBytes(byte(short >> 8), byte(short))
}

// Writes String in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteString(message string) {
	miniBuffer.Write([]byte(message))
}

// Writes a Int in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteInt(value int) {
	miniBuffer.WriteBytes(byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value))
}

// Writes Long in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteLong(long int64)  {
	miniBuffer.WriteBytes(byte(long >> 56), byte(long >> 48),
		byte(long >> 40), byte(long >> 32), byte(long >> 24),
		byte(long >> 16), byte(long >> 8), byte(long))
}

// To implement io.Writer
// Writes length of bytes, then the bytes.
func (miniBuffer *MiniBuffer) Write(bytes []byte) (n int, err error) {
	miniBuffer.WriteVarInt(len(bytes))
	miniBuffer.WriteBytes(bytes...)
	return len(bytes), nil
}

// Writes VarInt in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) WriteVarInt(value int) {

	for value & 0x80 != 0 {
		miniBuffer.WriteBytes(byte(value & 0x7F | 0x80))
		value >>= 7
	}

	miniBuffer.WriteBytes(byte(value & 0x7F))
}


// Reads Signed Short (int16) from the Buffer.
func (miniBuffer *MiniBuffer) ReadShort() int16 {
	return int16(miniBuffer.ReadUnsignedShort())
}

// Reads Unsigned Short (uint16) from the Buffer.
func (miniBuffer *MiniBuffer) ReadUnsignedShort() uint16 {
	return uint16(miniBuffer.ReadNext()) << 8 | uint16(miniBuffer.ReadNext())
}

// Reads Signed Long (int64) from the Buffer.
func (miniBuffer *MiniBuffer) ReadLong() int64 {
	return int64(miniBuffer.ReadUnsignedLong())
}

// TODO: Test
// Reads a Int in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) ReadInt() int {
	return int(miniBuffer.ReadNext() << 24 | miniBuffer.ReadNext() << 16 | miniBuffer.ReadNext() << 8 | miniBuffer.ReadNext())
}

// Reads Signed Long (int64) from the Buffer.
func (miniBuffer *MiniBuffer) ReadUnsignedLong() uint64 {
	return uint64(miniBuffer.ReadNext()) << 56 | uint64(miniBuffer.ReadNext()) << 48 |
		uint64(miniBuffer.ReadNext()) << 40 | uint64(miniBuffer.ReadNext()) << 32 |
		uint64(miniBuffer.ReadNext()) << 24 | uint64(miniBuffer.ReadNext()) << 16 |
		uint64(miniBuffer.ReadNext()) << 8 | uint64(miniBuffer.ReadNext())
}

// Reads String in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) ReadString() string {

	length := miniBuffer.ReadVarInt()

	readString := string(miniBuffer.Bytes[miniBuffer.ReaderIndex : miniBuffer.ReaderIndex + length])

	miniBuffer.ReaderIndex += length

	return readString
}

// Reads the next byte, clears and returns 0 if nothing is left.
func (miniBuffer *MiniBuffer) ReadNext() byte {

	if miniBuffer.ReaderIndex >= len(miniBuffer.Bytes) { miniBuffer.ClearAll(); return 0 }

	nextByte := miniBuffer.Bytes[miniBuffer.ReaderIndex]

	miniBuffer.ReaderIndex++

	if miniBuffer.WriterIndex > 0 { miniBuffer.WriterIndex-- }

	return nextByte
}

// Reads VarInt in respect of Minecraft's protocol.
func (miniBuffer *MiniBuffer) ReadVarInt() int {

	var result int

	for numRead := 0; ; {

		read := int(miniBuffer.ReadNext())

		result |= (read & 127) << uint(7 * numRead)

		if numRead++; numRead > 5 { println("VarInt is too big"); break }
		if read & 128 == 0 { break }

	}

	return result
}