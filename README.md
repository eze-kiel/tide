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
      - [Commands](#commands)
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

```
$ tide [filename]
```

### Shortcuts

#### Visual mode

|           Shortcut            | Action                                       |
| :---------------------------: | :------------------------------------------- |
|         <kbd>:</kbd>          | Open the command menu                        |
|         <kbd>I</kbd>          | Start inserting (switch to Edit mode)        |
|         <kbd>H</kbd>          | Move the cursor to the beginning of the line |
|         <kbd>L</kbd>          | Move the cursor to the end of the line       |
|         <kbd>T</kbd>          | Move the cursor to the top of the file       |
|         <kbd>E</kbd>          | Move the cursor to the end of the file       |
|         <kbd>O</kbd>          | Insert a new line under the cursor           |
| <kbd>Shift</kbd>+<kbd>O</kbd> | Insert a new line above the cursor           |
| <kbd>Ctrl</kbd>+<kbd>D</kbd>  | Fast jump downward                           |
| <kbd>Ctrl</kbd>+<kbd>U</kbd>  | Fast jump upward                             |
|         <kbd>D</kbd>          | Delete under cursor                          |

#### Edit mode

|    Shortcut    | Action                                          |
| :------------: | :---------------------------------------------- |
| <kbd>Esc</kbd> | Switch to Visual Mode (and autosave if enabled) |

#### Commands

|          Command           | Action                                    |
| :------------------------: | :---------------------------------------- |
|        `q`, `quit`         | Quit the editor                           |
|       `q!`, `quit!`        | Force quit the editor                     |
| `w [file]`, `write [file]` | Write changes to file                     |
|  `wq [file]`, `x [file]`   | Write changes to file and quit the editor |
|            `ge`            | Go to the end of the file                 |
|            `gt`            | Go to the top of the file                 |

## License

MIT