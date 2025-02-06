package main

import (
	"fmt"
	"os"

	"github.com/IkeVoodoo/go-copy/copy"
	"github.com/integrii/flaggy"
	cp "github.com/otiai10/copy"
)

func validateFilePath(path string, shouldExist bool) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if !shouldExist {
				return nil, nil
			}

			return nil, fmt.Errorf("path does not exist: %s", path)
		}

		return nil, fmt.Errorf("unable to access path: %s", path)
	}

	return info, nil
}

func exitOnError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type CopyOptions struct {
	overwriteExistingFiles bool
}

func handleCopy(source copy.ElementInfo, dest copy.ElementInfo, opts CopyOptions) error {
	if dest.Info != nil && !opts.overwriteExistingFiles {
		return fmt.Errorf("destination file exists, but overwriteExistingFiles is set to false. (Did you forget --overwrite-existing-files): %v", dest.Path)
	}

	err := os.RemoveAll(dest.Path)
	if err != nil {
		return err
	}

	return cp.Copy(source.Path, dest.Path, cp.Options{})
}

func main() {
	flaggy.SetName("copy")
	flaggy.SetName("Simple and intuitive CLI copy alternative")
	flaggy.SetVersion("1.0.0")

	var sourcePath string
	flaggy.AddPositionalValue(&sourcePath, "source_path", 1, true, "The source file or directory to copy from")

	var destPath string
	flaggy.AddPositionalValue(&destPath, "dest_path", 2, true, "The destination file or directory to copy to")

	var copyOptions = CopyOptions{}
	flaggy.Bool(&copyOptions.overwriteExistingFiles, "o", "overwrite-existing-files", "Should existing files have their contents overwritten? (Default: false)")

	flaggy.Parse()

	sourceInfo, err := validateFilePath(sourcePath, true)
	exitOnError(err)

	sourceElement := copy.ElementInfo{
		Path: sourcePath,
		Info: sourceInfo,
	}

	destInfo, err := validateFilePath(destPath, false)
	exitOnError(err)

	destElement := copy.ElementInfo{
		Path: destPath,
		Info: destInfo,
	}

	err = handleCopy(sourceElement, destElement, copyOptions)
	exitOnError(err)
}
