package converter

import (
	"fileconverter/internal/converterror"
	"fileconverter/internal/jsonfile"
	"fileconverter/internal/model"
	"fileconverter/internal/xlsxfile"
	"fmt"
	"log"
	"os"
)

type srcConverter interface {
	ToObjects() (model.Objects, error)
}

type dstConverter interface {
	ToFile(objects model.Objects) error
}

func Run() error {
	// распарсить аргументы командной строки
	src, dst, err := parseArgs()
	if err != nil {
		return fmt.Errorf("parseArgs: %v", err)
	}

	// открыть файл-источник
	err = src.openOSFile(os.O_RDONLY)
	if err != nil {
		return fmt.Errorf("src openOSFile: %v", err)
	}
	defer src.closeOSFile()

	// открыть файл-приемник
	err = dst.openOSFile(os.O_CREATE | os.O_RDWR)
	if err != nil {
		return fmt.Errorf("dst openOSFile: %v", err)
	}
	defer dst.closeOSFile()

	// определение правила конвертации
	from, to, err := convertationRule(src, dst)
	if err != nil {
		return fmt.Errorf("convertationRule: %v", err)
	}

	// конвертация
	err = convert(from, to)
	if err != nil {
		return fmt.Errorf("convert: %v", err)
	}

	log.Println("successful conversion")

	return nil
}

func convertationRule(src file, dst file) (srcConverter, dstConverter, error) {
	switch {
	case src.ext == jsonExt && dst.ext == xlsxExt:
		return jsonfile.New(src.osFile), xlsxfile.New(dst.osFile), nil
	}

	return nil, nil, fmt.Errorf("%w: %s -> %s", converterror.ErrNotImplementedRule, src.ext, dst.ext)
}

func convert(src srcConverter, dst dstConverter) error {
	objects, err := src.ToObjects()
	if err != nil {
		return fmt.Errorf("src ToObjects: %v", err)
	}
	err = dst.ToFile(objects)
	if err != nil {
		return fmt.Errorf("dst ToFile: %v", err)
	}

	return nil
}
