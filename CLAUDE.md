# CLAUDE.md

This file provides context for Claude Code when working on this project.

## Project Overview

Go AST Viewer is a GUI application that visualizes Go Abstract Syntax Trees. Users input Go code in txtar format on the left panel, and the parsed AST is displayed as a tree on the right panel.

## Tech Stack

- **GUI Framework**: [guigui](https://github.com/guigui-gui/guigui) - Pure Go immediate-mode GUI
- **AST Parsing**: Go standard library (`go/ast`, `go/parser`, `go/token`)
- **Input Format**: txtar (`golang.org/x/tools/txtar`)

## File Structure

```
goastviewer/
├── main.go          # Entry point, Root widget, app initialization
├── parser.go        # AST parsing logic, converts Go code to tree nodes
├── left_panel.go    # Left panel with text input for txtar code
├── right_panel.go   # Right panel with AST tree display
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md
```

## Key Components

### main.go
- `Root` widget: Main container with horizontal split layout
- Connects left panel's source changes to right panel's AST display

### parser.go
- `ParseTxtar()`: Parses txtar content into AST nodes
- `ASTNode`: Tree node structure for display
- `FlattenNodes()`: Converts tree to flat list for List widget

### left_panel.go
- `LeftPanel`: Contains text input for txtar format Go code
- `SetOnSourceChanged()`: Callback when user triggers AST parsing

### right_panel.go
- `RightPanel`: Displays parsed AST as tree using guigui's List widget
- `SetSource()`: Receives source code and triggers parsing

## guigui Patterns

- Widgets embed `guigui.DefaultWidget`
- `Build()`: Add child widgets, configure properties
- `Layout()`: Position child widgets within bounds
- Use `basicwidget.TextInput` for text editing
- Use `basicwidget.List` for tree/list display with `IndentLevel` for hierarchy

## Build & Run

```bash
go build
./goastviewer
```

## Common Tasks

### Adding new AST node types
Edit `parser.go`, add cases in `declToNode()`, `stmtToNode()`, or `exprToNode()`.

### Modifying UI layout
Edit `Layout()` methods in respective panel files. Use `guigui.LinearLayout` for simple layouts or manual `image.Rectangle` positioning.

### Changing default sample code
Edit `defaultSource` constant in `left_panel.go`.
