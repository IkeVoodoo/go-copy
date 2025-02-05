package main

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
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

func main() {
	flaggy.SetName("copy")
	flaggy.SetName("Simple and intuitive CLI copy alternative")
	flaggy.SetVersion("1.0.0")

	var sourcePath string
	flaggy.AddPositionalValue(&sourcePath, "source_path", 1, true, "The source file or directory to copy from")

	var destPath string
	flaggy.AddPositionalValue(&destPath, "dest_path", 2, true, "The destination file or directory to copy to")

	flaggy.Parse()

	sourceInfo, err := validateFilePath(sourcePath, true)
	exitOnError(err)

	destInfo, err := validateFilePath(destPath, false)
	exitOnError(err)

	fmt.Println(sourceInfo.Name())
	if destInfo == nil {
		fmt.Println(destPath)
	} else {
		fmt.Println(destInfo.Name())
	}
}
