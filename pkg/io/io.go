package io

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync/atomic"

	"github.com/pkg/errors"
)

const (
	bufferSize    = 1024
	maxPacketSize = 2<<15 - 1
)

type ReadWriter struct {
	r *bufio.Reader
	w *bufio.Writer

	enableCountSize bool
	readSizeCount   uint64
	writeSizeCount  uint64

	rcChan chan uint64
	wcChan chan uint64
}

func NewReadWriter(rw io.ReadWriter, enableCountSize bool) *ReadWriter {
	var readWriter = &ReadWriter{
		r:               bufio.NewReaderSize(rw, bufferSize),
		w:               bufio.NewWriterSize(rw, bufferSize),
		enableCountSize: enableCountSize,
	}

	if enableCountSize {
		readWriter.rcChan = make(chan uint64, bufferSize)
		readWriter.wcChan = make(chan uint64, bufferSize)
		go readWriter.count()
	}

	return readWriter
}

func (r *ReadWriter) count() {
	for {
		select {
		case n, ok := <-r.rcChan:
			if !ok {
				return
			}
			atomic.AddUint64(&r.readSizeCount, n)

		case n, ok := <-r.wcChan:
			if !ok {
				return
			}
			atomic.AddUint64(&r.writeSizeCount, n)
		}
	}
}

func (r *ReadWriter) addReadCount(n uint64) {
	if r.enableCountSize {
		r.rcChan <- n
	}
}

func (r *ReadWriter) addWriteCount(n uint64) {
	if r.enableCountSize {
		r.wcChan <- n
	}
}

// 头部2 个字节标志长度
func (r *ReadWriter) ReadOnePacket() ([]byte, error) {
	var head = make([]byte, 2)
	if _, err := io.ReadFull(r.r, head); err != nil {
		return nil, errors.Wrap(err, "read one packet header failed")
	}
	r.addReadCount(2)
	var length = int(binary.BigEndian.Uint16(head))
	var data = make([]byte, length)
	if n, err := io.ReadFull(r.r, data); err != nil {
		if n != 0 {
			r.addReadCount(uint64(n))
		}
		return nil, errors.Wrap(err, "read one packet body failed")
	}
	r.addReadCount(uint64(length))
	return data, nil
}

// TODO:: if client send the packet's length == maxPacketSize,
// it will be blocked. or will receive next package into one package.
func (r *ReadWriter) ReadPacket() ([]byte, error) {
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

func (r *ReadWriter) WritePacket(data []byte) error {
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
		r.addWriteCount(2)
		if _, err := r.w.Write(writeData); err != nil {
			return errors.Wrap(err, "failed to write body")
		}
		r.addWriteCount(uint64(len(writeData)))
	}
	return nil
}

func (r *ReadWriter) Flush() error {
	return r.w.Flush()
}
