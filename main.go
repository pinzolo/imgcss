package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var dir string
var file string
var suffix string
var validExts = []string{".png", ".jpg", ".jpeg"}

func main() {
	flag.StringVar(&file, "file", "", "directory")
	flag.StringVar(&file, "f", "", "directory")
	flag.StringVar(&dir, "dir", "", "directory")
	flag.StringVar(&dir, "d", "", "directory")
	flag.StringVar(&suffix, "suffix", "_image", "suffix")
	flag.Parse()

	err := proc()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func proc() error {
	var css string
	if file != "" {
		info, err := os.Stat(file)
		if err != nil {
			return err
		}
		css, err = cssFile(file, info)
		if err != nil {
			return err
		}
	} else if dir != "" {
		var list []string
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !isImageFile(path) {
				return nil
			}

			s, err := cssFile(path, info)
			if err != nil {
				return errors.Wrap(err, info.Name())
			}
			list = append(list, s)
			return nil
		})
		css = strings.Join(list, "\n")
	} else {
		return errors.New("required file or directory path")
	}

	fmt.Println(css)
	return nil
}

func cssFile(path string, info os.FileInfo) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return "", err
	}

	rect := img.Bounds()
	return fmt.Sprintf(".%s {\n  width: %dpx;\n  height: %dpx;\n}\n", className(info), rect.Dx(), rect.Dy()), nil
}

func isImageFile(fp string) bool {
	ext := strings.ToLower(filepath.Ext(fp))
	for _, ve := range validExts {
		if ve == ext {
			return true
		}
	}
	return false
}

func decoder(fp string) func(r io.Reader) (image.Image, error) {
	if strings.ToLower(filepath.Ext(fp)) == ".png" {
		return png.Decode
	}
	return jpeg.Decode
}

func className(fi os.FileInfo) string {
	return fi.Name()[0:len(fi.Name())-len(filepath.Ext(fi.Name()))] + suffix
}
