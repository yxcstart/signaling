package xrpc

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	HEADER_SIZE     = 36
	HEADER_MAGICNUM = 0xfb202309
)

type Header struct {
	Id       uint16
	Version  uint16
	LogId    uint32
	Provider []byte
	MagicNum uint32
	Reserved uint32
	Bodylen  uint32
}

func (h *Header) Marshal(b []byte) error {
	if len(b) < HEADER_SIZE {
		return fmt.Errorf("no enough buffer for header")
	}

	binary.LittleEndian.PutUint16(b[0:2], h.Id)
	binary.LittleEndian.PutUint16(b[2:4], h.Version)
	binary.LittleEndian.PutUint32(b[4:8], h.LogId)
	copy(b[8:24], h.Provider[:])
	binary.LittleEndian.PutUint32(b[24:28], h.MagicNum)
	binary.LittleEndian.PutUint32(b[28:32], h.Reserved)
	binary.LittleEndian.PutUint32(b[32:36], h.Bodylen)

	return nil
}

func (h *Header) Write(w io.Writer) (n int, err error) {
	var buf [HEADER_SIZE]byte
	if err = h.Marshal(buf[:]); err != nil {
		return 0, err
	}
	return w.Write(buf[:])
}
