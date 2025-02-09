package copy

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/integrii/flaggy"
)

type CopyOptions struct {
	OverwriteExistingFiles bool
	ScanSourcePath         bool
	ProgressBarVisible     bool
	BufferSize             int64
}

type CopyCallback func(path string, isDir bool, writtenBytes int) error

func (options *CopyOptions) InitializeOptionFlags() {
	flaggy.Bool(&options.OverwriteExistingFiles, "o", "overwrite-existing-files", "Should existing files have their contents overwritten?")
	flaggy.Bool(&options.ScanSourcePath, "s", "scan-source-path", "Should the source path be scanned to provide a progress bar and ETA?")
	flaggy.Bool(
		&options.ProgressBarVisible,
		"p",
		"progress-bar-visible",
		"Should the progress bar be visible? Only shows copied file count if --scan-source-path is missing.",
	)

	options.BufferSize = 8192
	flaggy.Int64(
		&options.BufferSize,
		"c",
		"chunk-size",
		"What chunk size (in bytes) should be used to copy files? A larger chunk size may result in faster speeds, at the cost of memory usage.",
	)
}

func copyFileUnsafe(source ElementInfo, dest ElementInfo, otherWriter io.Writer) error {
	sourceStream, err := os.Open(source.Path)
	if err != nil {
		return err
	}
	defer sourceStream.Close()

	destStream, err := os.Create(dest.Path)
	if err != nil {
		return err
	}
	defer destStream.Close()

	var writer io.Writer = destStream
	if otherWriter != nil {
		writer = io.MultiWriter(destStream, otherWriter)
	}

	_, err = io.Copy(writer, sourceStream)
	return err
}

func copyFile(source ElementInfo, dest ElementInfo, opts CopyOptions, otherWriter io.Writer) error {
	if dest.Info != nil && !opts.OverwriteExistingFiles {
		return fmt.Errorf("--overwrite-existing-files is missing, but dest exists: %v", dest.Path)
	}

	return copyFileUnsafe(source, dest, otherWriter)
}

func copyDir(source ElementInfo, dest ElementInfo, opts CopyOptions, otherWriter io.Writer) error {
	if dest.Info != nil && !opts.OverwriteExistingFiles {
		return fmt.Errorf("--overwrite-existing-files is missing, but dest exists: %v", dest.Path)
	}

	return filepath.WalkDir(source.Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		relative, err := filepath.Rel(source.Path, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dest.Path, relative)
		if d.IsDir() {
			err := os.Mkdir(target, info.Mode())
			if err != nil {
				return err
			}

			return nil
		}

		return copyFileUnsafe(ElementInfo{
			Path: path,
			Info: info,
		}, ElementInfo{
			Path: target,
			Info: nil,
		}, otherWriter)
	})
}

func Copy(source ElementInfo, dest ElementInfo, opts CopyOptions, otherWriter io.Writer) error {
	if source.Info.IsDir() {
		return copyDir(source, dest, opts, otherWriter)
	}

	return copyFile(source, dest, opts, otherWriter)
}
