package copy

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type DiscoverResult struct {
	FileTotal      uint32
	DirectoryTotal uint32
	SizeTotal      uint64
}

const (
	BYTE     = 1.0
	KILOBYTE = BYTE * 1000
	MEGABYTE = KILOBYTE * 1000
	GIGABYTE = MEGABYTE * 1000
	TERABYTE = GIGABYTE * 1000
)

var messagePrinter = message.NewPrinter(language.English)

func (r *DiscoverResult) PrettySizeTotal() string {
	if r.SizeTotal == 0 {
		return "0 Bytes"
	}

	unit := ""
	value := float32(r.SizeTotal)

	switch {
	case r.SizeTotal >= TERABYTE:
		unit = "TB"
		value /= float32(TERABYTE)
	case r.SizeTotal >= GIGABYTE:
		unit = "GB"
		value /= float32(GIGABYTE)
	case r.SizeTotal >= MEGABYTE:
		unit = "MB"
		value /= float32(MEGABYTE)
	case r.SizeTotal >= KILOBYTE:
		unit = "KB"
		value /= float32(KILOBYTE)
	case r.SizeTotal == 1:
		unit = "Byte"
	case r.SizeTotal > 0:
		unit = "Bytes"
	}

	stringValue := strings.TrimSuffix(
		messagePrinter.Sprintf("%.2f", value), ".00",
	)

	return fmt.Sprintf("%s %s", stringValue, unit)
}

func (r *DiscoverResult) PrettyPathTotal() string {
	return messagePrinter.Sprintf("%.2d", r.FileTotal+r.DirectoryTotal)
}

func DiscoverSource(source ElementInfo) (*DiscoverResult, error) {
	var result DiscoverResult
	if !source.Info.IsDir() {
		result.FileTotal = 1
		result.DirectoryTotal = 0
		result.SizeTotal = uint64(source.Info.Size())
		return &result, nil
	}

	err := filepath.WalkDir(source.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			result.DirectoryTotal++
			return nil
		}

		result.FileTotal++

		info, err := d.Info()
		if err != nil {
			return err
		}
		result.SizeTotal += uint64(info.Size())

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}
