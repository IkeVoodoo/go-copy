package main

import (
	"fmt"
	"os"

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

type ElementInfo struct {
	path string
	info os.FileInfo
}

func handleCopy(source ElementInfo, dest ElementInfo, opts CopyOptions) error {
	if dest.info != nil && !opts.overwriteExistingFiles {
		return fmt.Errorf("destination file exists, but overwriteExistingFiles is set to false. (Did you forget --overwrite-existing-files): %v", dest.path)
	}

	err := os.RemoveAll(dest.path)
	if err != nil {
		return err
	}

	return cp.Copy(source.path, dest.path, cp.Options{})
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

	sourceElement := ElementInfo{
		path: sourcePath,
		info: sourceInfo,
	}

	destInfo, err := validateFilePath(destPath, false)
	exitOnError(err)

	destElement := ElementInfo{
		path: destPath,
		info: destInfo,
	}

	err = handleCopy(sourceElement, destElement, copyOptions)
	exitOnError(err)
}
