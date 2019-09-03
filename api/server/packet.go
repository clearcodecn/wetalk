package server

import (
	"bufio"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

const (
	bufferSize    = 1024
	maxPacketSize = 2<<15 - 1
)

type packet struct {
	r *bufio.Reader
	w *bufio.Writer
	rcChan chan uint64
	wcChan chan uint64
}

func NewPacket(rw io.ReadWriter) *packet {
	var readWriter = &packet{
		r:               bufio.NewReaderSize(rw, bufferSize),
		w:               bufio.NewWriterSize(rw, bufferSize),
	}
	return readWriter
}

// 头部2 个字节标志长度
func (r *packet) ReadOnePacket() ([]byte, error) {
	var head = make([]byte, 2)
	if _, err := io.ReadFull(r.r, head); err != nil {
		return nil, errors.Wrap(err, "read one packet header failed")
	}
	var length = int(binary.BigEndian.Uint16(head))
	var data = make([]byte, length)
	if _, err := io.ReadFull(r.r, data); err != nil {
		return nil, errors.Wrap(err, "read one packet body failed")
	}
	return data, nil
}

// TODO:: if client send the packet's length == maxPacketSize,
// it will be blocked. or will receive next package into one package.
func (r *packet) ReadPacket() ([]byte, error) {
	packet, err := r.ReadOnePacket()
	if err != nil {
		return nil, err
	}
	if len(packet) < maxPacketSize {
		return packet, nil
	}
	for {
		data, err := r.ReadOnePacket()
		if err != nil {
			return nil, err
		}
		packet = append(packet, data...)
		if len(data) < maxPacketSize {
			break
		}
	}
	return packet, nil
}

func (r *packet) WritePacket(data []byte) error {
	var (
		head       = make([]byte, 2)
		length     = len(data)
		leftLength = length
	)
	for leftLength != 0 {
		var writeData []byte
		if leftLength >= maxPacketSize {
			leftLength -= maxPacketSize
			length = maxPacketSize
			writeData = data[:maxPacketSize]
			data = data[maxPacketSize:]
		} else {
			length = leftLength
			leftLength = 0
			writeData = data
		}
		binary.BigEndian.PutUint16(head, uint16(length))
		if _, err := r.w.Write(head); err != nil {
			return errors.Wrap(err, "failed to write header")
		}
		if _, err := r.w.Write(writeData); err != nil {
			return errors.Wrap(err, "failed to write body")
		}
	}
	return nil
}

func (r *packet) Flush() error {
	return r.w.Flush()
}
