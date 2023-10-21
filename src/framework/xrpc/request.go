package xrpc

import (
	"bytes"
	"io"
)

type Request struct {
	Header Header
	Body   io.Reader
}

func NewRequest(body io.Reader, logId uint32) *Request {
	req := &Request{}
	req.Header.LogId = logId
	req.Header.MagicNum = HEADER_MAGICNUM

	if body != nil {
		switch v := body.(type) {
		case *bytes.Reader:
			req.Header.Bodylen = uint32(v.Len())
		default:
			return nil
		}

		req.Body = io.LimitReader(body, int64(req.Header.Bodylen))
	}

	return req
}

func (r *Request) Write(w io.Writer) (n int, err error) {
	if n, err = r.Header.Write(w); err != nil {
		return 0, err
	}
	written, err := io.Copy(w, r.Body)
	return int(written), err
}
