# FileBot

FileBot is simple file manager to automation your boring action with files e.g. moving to trash or another director

# Configuration

## Flags

```text
Usage:
  filebot [flags]

Flags:
  -c, --config string             specific path for filebot configuration (default "/home/wittano/.config/filebot/setting.toml")
  -h, --help                      help for filebot
  -l, --log string                path to log file
      --logLevel string           log level (default "INFO")
  -u, --updateInterval duration   set time after filebot should be refresh watched file state (default 10m0s)
```

## Configuration file

```toml
# Example configuration for moving some files to another location
[Example]
src = ["/path/to/example.txt"]
dest = "/path/to/destination/directory"

# You can use standard regular expression seeking your files
[SearchFilesByRegex]
src = ["/path/to/example.*", "/path/to/[0-9]ystem.txt"]
dest = "/path/to/destination/directory"

# Automaticlly move your unnecessary files to Trash directory
[MoveToTrash]
src = ["/path/to/your/junk-file.*"]
moveToTrash = true
after = 32 # Time a last modification(in days), after that files will be moved to trash

# Some files you don't want to move or touch by program. You can add specified file to exception list
[Exception]
src = ["/some/files.*"]
dest = "/path/to/desctination"
exceptions = ["files.pdf", "files.ini"] # Put only file name
```