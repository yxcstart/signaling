package xrpc

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
