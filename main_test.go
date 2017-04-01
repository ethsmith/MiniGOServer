package main_test

import (
	"testing"
	"bytes"
	"compress/zlib"
	"sync"
	"math"
	"math/rand"
	"GoMiniServer/networking"
	"crypto/sha1"
	"GoMiniServerRemake/src/server"
	"github.com/jinzhu/gorm"
)

// Compresses a byte array of characters for testing compression with Zlib
func TestZlibCompress(testing *testing.T) {
	var blankBuffer bytes.Buffer
	w := zlib.NewWriter(&blankBuffer)
	w.Write([]byte("This shalt be compressed."))
	w.Close()
}

// Tests encryption of the word "Hi"
func TestEncryption(testing *testing.T) {

	encrypt := "Hi"
	sha := sha1.New()

	sha.Write([]byte(encrypt))

	//hash := sha.Sum(nil)

}

// Tries to start the server, and store it's data in GORM..
func TestServerStart(testing *testing.T) {

	server := networking.NewMiniServer("127.0.0.1", 25565, "")

	DB, _ := gorm.Open("sqlite3", "MiniServer.db")

	DB.Create(server)

	DB.Close()

	if !server.Start() {
		testing.Error("Error starting TCP server.")

	} else { testing.Log("Works?") }

	server.Stop()

	// TODO: Fix Data Race

}

// Tries to get the absolute value of a set of values.
func TestToAbsolute(testing *testing.T) {
	var randomNumber = rand.Intn(9)

	if randomNumber == 0 { testing.Error("0 has been found D:"); return }

	a := []float64{ 0, -1, -3.00000 }

	testing.Log(FloatToAbsolute1(a))
	testing.Log(FloatToAbsolute2(a))

	var converted = FloatToAbsolute3(a)
	for val := range converted { testing.Log(val) }
}

// Tries to get the absolute value of a set of values.
func FloatToAbsolute1(floats []float64) []float64 {
	var unsignedFloats = make([]float64, len(floats))
	for i, v := range floats { unsignedFloats[i] = math.Abs(v) }
	return unsignedFloats
}

// Tries to get the absolute value of a set of values.
func FloatToAbsolute2(floats []float64) []float64 {

	var size = len(floats)
	var unsignedFloats = make([]float64, size)

	wg := sync.WaitGroup{}

	wg.Add(size)

	calc := func(index int, val float64){
		unsignedFloats[index] = math.Abs(val)
		wg.Done()
	}

	for i, v := range floats { go calc(i, v) }

	wg.Wait()
	return unsignedFloats
}

// Tries to get the absolute value of a set of values.
func FloatToAbsolute3(floats []float64) chan float64 {
	var returns = make(chan float64, len(floats))

	go func() {
		for _, v := range floats { returns <- math.Abs(v) }
		close(returns)
	}()

	return returns
}


