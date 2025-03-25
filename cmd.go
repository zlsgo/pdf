package pdf

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sohaha/zlsgo/zfile"
	"github.com/sohaha/zlsgo/zshell"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/zutil"
)

var pdftotext = zutil.Once(func() func(path string) (string, error) {
	// apt install poppler-utils

	code, out, errout, err := zshell.Run("pdftotext -v")
	if err != nil || code != 0 || !(strings.Contains(out, "poppler") || strings.Contains(errout, "poppler")) {
		return nil
	}

	return func(path string) (string, error) {
		code, out, errout, err := zshell.Run(`pdftotext -nopgbrk -raw -enc UTF-8 "` + path + `" -`)
		if code == 99 {
			code, out, errout, err = zshell.Run(`pdftotext -nopgbrk -raw -enc GBK "` + path + `" -`)
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

var pdftoppm = zutil.Once(func() func(path string) (map[string]string, error) {
	code, out, errout, err := zshell.Run("pdftoppm -v")
	if err != nil || code != 0 || !(strings.Contains(out, "pdftoppm") || strings.Contains(errout, "pdftoppm")) {
		return nil
	}

	tmp := zfile.TmpPath()
	return func(path string) (map[string]string, error) {
		dir := zfile.RealPathMkdir(tmp+"/"+zstring.Md5(path), true)
		code, _, errout, err := zshell.Run(`pdftoppm -png "` + path + `" ` + dir)
		if err != nil {
			return map[string]string{}, err
		}

		if code != 0 {
			return map[string]string{}, errors.New(errout)
		}

		defer func() {
			_ = zfile.Rmdir(dir)
		}()

		images := make(map[string]string, 0)
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			base := filepath.Base(path)
			ext := filepath.Ext(base)
			if strings.ToLower(ext) != ".png" {
				return nil
			}

			index := strings.TrimPrefix(strings.TrimSuffix(base, ext), "-")
			base64, _ := zstring.Img2Base64(path)
			images[index] = base64
			return nil
		})

		return images, nil
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
