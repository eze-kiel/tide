# tide

_my own minimalist text editor_

> [!CAUTION]
> This project is in early development and probably really buggy.

## Usage

### Build

- As there's no CI yet, you need to build the binary manually. `go` is required.

```
$ go build -o build/ed .
```

- Then you can copy the bin somewhere in your path:

```
$ cp build/tide /<in>/<your>/<$PATH>/tide
```

### Run

For now, you must provide a filename in arg 1:

```
$ tide <filename>
```

### Shortcuts

- Ctrl + W: Save
- Ctrl + Q (or Esc): Quit
- Ctrl + L: Go to the end of the line
- Ctrl + H: Go to the beginning of the line

## License

MIT