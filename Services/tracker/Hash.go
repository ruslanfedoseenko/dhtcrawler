package tracker

import (
	"encoding/hex"
	"fmt"
)

type Hash [20]byte

func (h Hash) Bytes() []byte {
	return h[:]
}

func (h Hash) AsString() string {
	return string(h[:])
}

func (h Hash) HexString() string {
	return fmt.Sprintf("%x", h[:])
}

func (h *Hash) FromString(s string) (err error) {
	if len(s) != 20 {
		err = fmt.Errorf("string has bad length: %d", len(s))
		return
	}
	for i := 0; i < len(s); i++ {
		h[i] = s[i]
	}
	return
}

func (h *Hash) FromHexString(s string) (err error) {
	if len(s) != 40 {
		err = fmt.Errorf("hash hex string has bad length: %d", len(s))
		return
	}
	n, err := hex.Decode(h[:], []byte(s))
	if err != nil {
		return
	}
	if n != 20 {
		panic(n)
	}
	return
}
