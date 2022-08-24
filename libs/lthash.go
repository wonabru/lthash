package main // import "lukechampine.com/lthash"

import (
	"encoding/binary"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/blake2b"
)

const (
	BYTESIZE = 32
)

// A Hash is an instance of LtHash.
type Hash interface {
	// Add adds p to the checksum.
	Add(p []byte)
	// Remove removes p from the checksum.
	Remove(p []byte)
	// Sum appends the current checksum to b and returns it.
	Sum(b []byte) []byte
	// SetState sets the current checksum.
	SetState(state []byte)
}

type hash16 struct {
	sum  [BYTESIZE]byte
	hbuf [BYTESIZE]byte
	xof  blake2b.XOF
}

func (h *hash16) hashObject(p []byte) *[BYTESIZE]byte {
	h.xof.Reset()
	h.xof.Write(p)
	_, err := io.ReadFull(h.xof, h.hbuf[:])
	if err != nil {
		panic(err)
	}
	return &h.hbuf
}

// Add implements Hash.

func (h *hash16) Add(p []byte) {
	add16(&h.sum, h.hashObject(p))
}

// Remove implements Hash.

func (h *hash16) Remove(p []byte) {
	sub16(&h.sum, h.hashObject(p))
}

// Sum implements Hash.

func (h *hash16) Sum(b []byte) []byte {
	return append(b, h.sum[:]...)
}

// SetState implements Hash.
func (h *hash16) SetState(state []byte) {
	for i := range &h.sum {
		h.sum[i] = 0
	}
	copy(h.sum[:], state)
}

// New16 returns an instance of lthash16.

func New16() Hash {
	xof, _ := blake2b.NewXOF(BYTESIZE, nil)
	return &hash16{xof: xof}
}

func add16(x, y *[BYTESIZE]byte) {
	for i := 0; i < BYTESIZE; i += 2 {
		xi, yi := x[i:i+2], y[i:i+2]
		sum := binary.LittleEndian.Uint16(xi) + binary.LittleEndian.Uint16(yi)
		binary.LittleEndian.PutUint16(xi, sum)
	}
}

func sub16(x, y *[BYTESIZE]byte) {
	for i := 0; i < BYTESIZE; i += 2 {
		xi, yi := x[i:i+2], y[i:i+2]
		sum := binary.LittleEndian.Uint16(xi) - binary.LittleEndian.Uint16(yi)
		binary.LittleEndian.PutUint16(xi, sum)
	}
}

func main() {

	h := New16()

	h.Add([]byte("abc"))

	// for i := 0; i < 1000000; i += 1 {
	// 	h.Add([]byte(string(i)))
	// }
	s2 := h.Sum(nil)

	h2 := New16()
	h2.Add([]byte("abc"))
	h2.Add([]byte("abc"))
	h2.Remove([]byte("bca"))
	// for i := 1; i < 1000001; i += 1 {
	// 	h2.Add([]byte(string(i)))
	// }
	// h2.Remove([]byte(string(1000000)))
	// h2.Add([]byte(string(0)))
	s3 := h2.Sum(nil)

	println(hex.EncodeToString(s2))
	println(hex.EncodeToString(s3))

}
