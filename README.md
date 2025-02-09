# go-copy
Simple and intuitive CLI copy alternative

## Usage
```
go-copy - Simple and intuitive CLI copy alternative

  Usage:
    go-copy [source_path] [dest_path]

  Positional Variables: 
    source_path   The source file or directory to copy from (Required)
    dest_path     The destination file or directory to copy to (Required)

  Flags: 
       --version                    Displays the program version string.
    -h --help                       Displays help with available flag, subcommand, and positional value parameters.
    -o --overwrite-existing-files   Should existing files have their contents overwritten?
    -s --scan-source-path           Should the source path be scanned to provide a progress bar and ETA?
    -p --progress-bar-visible       Should the progress bar be visible? Only shows copied file count if --scan-source-path is missing.
    -c --chunk-size                 What chunk size (in bytes) should be used to copy files? A larger chunk size may result in faster speeds, at the cost of memory usage. (default: 8192)
    -a --all                        Should all flags be enabled?
```

## Installing
Installing `go-copy` is pretty straightforward:

### Install with go
```sh
go install github.com/IkeVoodoo/go-copy/
```

### Install from source
```sh
git clone https://github.com/IkeVoodoo/go-copy/
cd co-copy
go install
```

### From releases
<TODO: Releases not published yet!>

### NixOS
<TODO: flake not ready yet!>


## Building
To build `go-copy` from source do:

```sh
git clone https://github.com/IkeVoodoo/go-copy/
cd co-copy
go build
```
