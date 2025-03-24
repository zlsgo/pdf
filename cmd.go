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
		code, out, errout, err := zshell.Run(`pdftotext -nopgbrk -enc UTF-8 "` + path + `" -`)
		if code == 99 {
			code, out, errout, err = zshell.Run(`pdftotext -nopgbrk -enc GBK "` + path + `" -`)
		}
		if err != nil {
			return "", err
		}

		if code != 0 {
			return "", errors.New(errout)
		}
		return out, nil
	}
})

var markitdown = zutil.Once(func() func(path string) (string, error) {
	code, out, _, err := zshell.Run("markitdown -v")
	if err != nil || code != 0 || !(strings.Contains(out, "markitdown")) {
		return nil
	}

	return func(path string) (string, error) {
		code, out, errout, err := zshell.Run("markitdown " + path)
		if err != nil {
			return "", err
		}
		if code != 0 {
			return "", errors.New(errout)
		}
		return out, nil
	}
})
