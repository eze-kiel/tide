# tide

_my own minimalist text editor_

> [!WARNING]
> This project is in early development and is probably really buggy.

- [tide](#tide)
  - [Usage](#usage)
    - [Build](#build)
    - [Run](#run)
    - [Shortcuts](#shortcuts)
      - [Visual mode](#visual-mode)
      - [Edit mode](#edit-mode)
  - [License](#license)

## Usage

### Build

- As there's no CI yet, you need to build the binary manually. `go` is required.

```
$ go build -o build/tide .
```

- Then you can copy the binary somewhere in your path:

```
$ cp ./build/tide /<in your $PATH>/tide
```

### Run

For now, you **must** provide a filename in first argument:

```
$ tide <filename>
```

### Shortcuts

#### Visual mode

|   Shortcut   | Action                                       |
| :----------: | :------------------------------------------- |
| <kbd>Q</kbd> | Quit the editor                              |
| <kbd>I</kbd> | Start inserting (switch to Edit mode)        |
| <kbd>L</kbd> | Move the cursor to the end of the line       |
| <kbd>H</kbd> | Move the cursor to the beginning of the line |
| <kbd>D</kbd> | Fast jump downward (default: 10 lines)       |
| <kbd>U</kbd> | Fast jump upward (default: 10 lines)         |
| <kbd>W</kbd> | Save buffer to file                          |

#### Edit mode

|           Shortcut           | Action                                      |
| :--------------------------: | :------------------------------------------ |
|        <kbd>Esc</kbd>        | Switch to Visual Mode (Autosave if enabled) |
| <kbd>Ctrl</kbd>+<kbd>Q</kbd> | Quit the editor                             |
| <kbd>Ctrl</kbd>+<kbd>X</kbd> | Delete current line                         |
| <kbd>Ctrl</kbd>+<kbd>L</kbd> | Move cursor to the end of the line          |
| <kbd>Ctrl</kbd>+<kbd>H</kbd> | Move cursor to the beginning of the line    |
| <kbd>Ctrl</kbd>+<kbd>D</kbd> | Fast jump downward (default: 10 lines)      |
| <kbd>Ctrl</kbd>+<kbd>U</kbd> | Fast jump upward (default: 10 lines)        |
| <kbd>Ctrl</kbd>+<kbd>W</kbd> | Save buffer to file                         |

## License

MIT