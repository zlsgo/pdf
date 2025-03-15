package pdf

import (
	"bytes"
	"errors"

	"github.com/sohaha/zlsgo/zfile"
)

func ToText(path string) (string, error) {
	path = zfile.RealPath(path)
	if !zfile.FileExist(path) {
		return "", errors.New("file not exists")
	}

	s := pdftotext()
	if s != nil {
		return s(path)
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
	return buf.String(), nil
}
