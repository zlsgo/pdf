package pdf

import (
	"bytes"
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zstring"
)

func ToText(path string, removeChars ...bool) (string, error) {
	path = zfile.RealPath(path)
	if !zfile.FileExist(path) {
		return "", errors.New("file not exists")
	}

	remove := len(removeChars) > 0 && removeChars[0]

	s := pdftotext()
	if s != nil {
		resp, err := s(path)
		if err == nil {
			return removeExtraChars(zstring.String2Bytes(resp), remove), nil
		}
	}

	s = markitdown()
	if s != nil {
		resp, err := s(path)
		if err == nil {
			return removeExtraChars(zstring.String2Bytes(resp), remove), nil
		}
	}

	f, r, err := Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	buf.ReadFrom(b)
	return removeExtraChars(buf.Bytes(), remove), nil
}

func ToImageBase64(path string) (map[string]string, error) {
	s := pdftoppm()
	if s == nil {
		return map[string]string{}, errors.New("please install poppler first")
	}

	return s(path)
}

func removeExtraChars(text []byte, remove bool) string {
	if !remove {
		return zstring.Bytes2String(text)
	}

	lines := bytes.Split(text, []byte("\n"))
	k := make(map[string]int, len(lines))
	nlines := make([]string, 0, len(lines))
	for i := range lines {
		line := strings.TrimSpace(zstring.Bytes2String(lines[i]))
		if len(line) == 0 {
			continue
		}

		_, ok := k[line]
		if !ok {
			k[line] = 0
		} else {
			k[line]++
		}
		if k[line] >= 2 {
			continue
		}

		nlines = append(nlines, line)
	}

	return strings.Join(nlines, "\n")
}
