package main

import (
	"fmt"
	"io"
	"os"

	"github.com/IkeVoodoo/go-copy/copy"
	"github.com/integrii/flaggy"
	"github.com/schollz/progressbar/v3"
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

func expandOptions(slice []string) []string {
	var result []string

	for _, item := range slice {
		if len(item) > 2 && item[0] == '-' && item[1] != '-' {
			for _, char := range item[1:] {
				result = append(result, "-"+string(char))
			}

			continue
		}

		result = append(result, item)
	}

	return result
}

func main() {
	flaggy.SetName("go-copy")
	flaggy.SetDescription("Simple and intuitive CLI copy alternative")
	flaggy.SetVersion("1.0.0")

	var sourcePath string
	flaggy.AddPositionalValue(&sourcePath, "source_path", 1, true, "The source file or directory to copy from")

	var destPath string
	flaggy.AddPositionalValue(&destPath, "dest_path", 2, true, "The destination file or directory to copy to")

	var copyOptions copy.CopyOptions
	copyOptions.InitializeOptionFlags()

	flaggy.ParseArgs(expandOptions(os.Args[1:]))

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

	if copyOptions.EnableAll {
		copyOptions.OverwriteExistingFiles = true
		copyOptions.ProgressBarVisible = true
		copyOptions.ScanSourcePath = true
	}

	max := int64(-1)
	if copyOptions.ScanSourcePath {
		var discovery copy.DiscoverResult
		discovery.PopulateFromPath(sourceElement)

		max = int64(discovery.SizeTotal)
	}

	var progress io.Writer
	if copyOptions.ProgressBarVisible {
		progress = progressbar.DefaultBytes(max, fmt.Sprintf("Copying %v to %v", sourcePath, destPath))
	}

	err = copy.Copy(sourceElement, destElement, copyOptions, progress)
	exitOnError(err)
}
