package pdf

import (
	"errors"
	"strings"

	"github.com/sohaha/zlsgo/zshell"
	"github.com/sohaha/zlsgo/zutil"
)

var pdftotext = zutil.Once(func() func(path string) (string, error) {
	// apt install poppler-utils

	code, out, errout, err := zshell.Run("pdftotext -v")
	if err != nil || code != 0 || !(strings.Contains(out, "poppler") || strings.Contains(errout, "poppler")) {
		return nil
	}

	return func(path string) (string, error) {
		code, out, errout, err := zshell.Run("pdftotext -enc UTF-8 " + path + " -")
		if err != nil {
			return "", err
		}
		if code != 0 {
			return "", errors.New(errout)
		}
		return out, nil
	}

})
