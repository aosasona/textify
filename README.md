# Textify

Take all the source code in a directory matching the provided parameters and put them all in a text file.

## Why?

Because Uni?

## Building

```sh
go build -o ./bin
```

## Usage

```sh
./bin [options] path/to/foo
```

### Flags

- `file` - output file, defaults `file.txt`
- `extension` - extensions to look for
- `ignore` - a list of directories to ignore separated by commas
