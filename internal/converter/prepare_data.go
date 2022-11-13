package converter

import (
	"fileconverter/internal/converterror"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	path   string
	ext    string
	osFile *os.File
}

func parseArgs() (src file, dst file, err error) {
	flag.Parse()
	args := flag.Args()

	lenArgs := len(args)

	if lenArgs != 2 {
		return file{}, file{}, fmt.Errorf("%w: %v", converterror.ErrInvalidArgsCount, lenArgs)
	}

	src, err = newFile(args[0])
	if err != nil {
		return file{}, file{}, fmt.Errorf("parseFileArg src %v", err)
	}

	dst, err = newFile(args[1])
	if err != nil {
		return file{}, file{}, fmt.Errorf("parseFileArg dst %v", err)
	}

	if src.ext == dst.ext {
		return file{}, file{}, converterror.ErrEqualExtentions
	}

	return src, dst, nil
}

func newFile(path string) (file, error) {
	ext := filepath.Ext(path)
	if !checkExtension(ext) {
		return file{}, fmt.Errorf("%w: %s", converterror.ErrNotAllowedExtention, ext)
	}

	return file{
		path:   path,
		ext:    ext,
		osFile: nil,
	}, nil
}

func checkExtension(ext string) bool {
	switch strings.ToLower(ext) {
	case jsonExt:
		return true
	case xlsxExt:
		return true
	default:
		return false
	}
}

func (f *file) openOSFile(flag int) (err error) {
	f.osFile, err = os.OpenFile(f.path, flag, 0666)
	if err != nil {
		return fmt.Errorf("os OpenFile: %v", err)
	}

	return nil
}

func (f *file) closeOSFile() {
	err := f.osFile.Close()
	if err != nil {
		log.Printf("defer close file error: %v", err)
	}
}
