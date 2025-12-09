# Go AST Viewer

A GUI tool for visualizing Go Abstract Syntax Trees (AST) from txtar format input.

## Features

- Split-pane interface with text editor on the left and AST tree on the right
- Parse Go code in txtar format
- Display AST as an interactive tree view
- Support for multiple Go files in a single txtar archive

## Requirements

- Go 1.21 or later
- macOS, Linux, or Windows

## Installation

```bash
git clone https://github.com/knsh14/goastviewer.git
cd goastviewer
go mod tidy
go build
```

Or install directly:

```bash
go install github.com/knsh14/goastviewer@latest
```

## Usage

```bash
./goastviewer
```

Or run directly:

```bash
go run .
```

## txtar Format

The tool accepts Go code in [txtar format](https://pkg.go.dev/golang.org/x/tools/txtar). Example:

```
-- main.go --
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    p := &Person{Name: "Alice", Age: 30}
    fmt.Println(p.Name)
}
```

You can include multiple files:

```
-- main.go --
package main

func main() {
    hello()
}

-- hello.go --
package main

import "fmt"

func hello() {
    fmt.Println("Hello!")
}
```

## AST Tree Display

The AST tree shows:

- File structure
- Package declarations
- Import statements
- Type definitions (struct, interface, etc.)
- Function and method declarations
- Statements (if, for, return, etc.)
- Expressions (function calls, operators, etc.)

## Dependencies

- [guigui](https://github.com/guigui-gui/guigui) - Pure Go GUI framework
- [ebiten](https://github.com/hajimehoshi/ebiten) - 2D game engine (used by guigui)
- Go standard library (`go/ast`, `go/parser`, `go/token`)
- [golang.org/x/tools/txtar](https://pkg.go.dev/golang.org/x/tools/txtar) - txtar format parser

## License

Apache-2.0
