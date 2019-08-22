package io

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewReadWriter(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	rw1 := NewReadWriter(buf, true)
	require.NotNil(t, rw1)
	var buf2 = bytes.NewBuffer(nil)
	rw2 := NewReadWriter(buf2, false)
	require.NotNil(t, rw2)
}

func TestReadWriter(t *testing.T) {
	var buf = bytes.NewBuffer(nil)
	rw1 := NewReadWriter(buf, true)
	require.NotNil(t, rw1)

	d1 := newRandomSizeData(80)
	err := rw1.WritePacket(d1)
	require.Nil(t, err)

	err = rw1.Flush()
	require.Nil(t, err)

	d2 := newRandomSizeData(1025)
	err = rw1.WritePacket(d2)
	require.Nil(t, err)

	err = rw1.Flush()
	require.Nil(t, err)

	d3 := newRandomSizeData(maxPacketSize*3 + 4)
	err = rw1.WritePacket(d3)
	require.Nil(t, err)

	err = rw1.Flush()
	require.Nil(t, err)

	d4 := newRandomSizeData(maxPacketSize*2 + 1024)
	err = rw1.WritePacket(d4)
	require.Nil(t, err)

	err = rw1.Flush()
	require.Nil(t, err)

	r1, err := rw1.ReadPacket()
	require.Nil(t, err)
	require.Equal(t, d1, r1)

	r2, err := rw1.ReadPacket()
	require.Nil(t, err)
	require.Equal(t, d2, r2)

	r3, err := rw1.ReadPacket()

	require.Nil(t, err)
	require.Equal(t, d3, r3)

	r4, err := rw1.ReadPacket()
	require.Nil(t, err)
	require.Equal(t, d4, r4)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
func newRandomSizeData(size int) []byte {
	var b = make([]byte, size)
	rand.Read(b)
	return b
}
