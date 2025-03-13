package pdf

import (
	"github.com/sohaha/zlsgo/zutil"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

type utf16hDecode struct{}

var utf16h = zutil.Once(func() *encoding.Decoder {
	return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
})

func (e *utf16hDecode) Decode(raw string) string {
	text, err := utf16h().String(raw)
	if err != nil {
		return raw
	}

	return text
}
