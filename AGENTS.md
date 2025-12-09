# AGENTS.md

Guidelines for AI coding agents working on this project.

## Project Context

This is a Go AST Viewer GUI application built with the guigui framework. It parses Go code in txtar format and displays the AST as an interactive tree.

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable names
- Keep functions focused and small
- Add comments for non-obvious logic

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                      Root                           │
│  ┌──────────────────┐  ┌─────────────────────────┐  │
│  │   LeftPanel      │  │     RightPanel          │  │
│  │                  │  │                         │  │
│  │  - Title         │  │  - Title                │  │
│  │  - TextInput     │──│  - Tree List (AST)      │  │
│  │  - Parse Button  │  │  - Error display        │  │
│  │                  │  │                         │  │
│  └──────────────────┘  └─────────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

## Key Files

| File | Purpose |
|------|---------|
| `main.go` | App entry point, Root widget |
| `parser.go` | txtar parsing, AST conversion |
| `left_panel.go` | Text input UI |
| `right_panel.go` | AST tree display UI |

## guigui Framework Patterns

### Widget Structure
```go
type MyWidget struct {
    guigui.DefaultWidget
    // child widgets
    // state
}

func (w *MyWidget) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
    adder.AddChild(&w.childWidget)
    // configure widgets
    return nil
}

func (w *MyWidget) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
    // position child widgets
    layouter.LayoutWidget(&w.childWidget, bounds)
}
```

### Common Widgets
- `basicwidget.Text` - Static text display
- `basicwidget.TextInput` - Editable text (use `SetMultiline(true)` for multi-line)
- `basicwidget.Button` - Clickable button
- `basicwidget.List[T]` - List/tree display with `IndentLevel` for hierarchy
- `basicwidget.Panel` - Scrollable container

### Layout
- `guigui.LinearLayout` for simple vertical/horizontal layouts
- Manual `image.Rectangle` positioning for complex layouts
- `basicwidget.UnitSize(context)` for consistent spacing

## Testing Changes

```bash
# Build
go build

# Run
./goastviewer

# Or combined
go run .
```

## Common Modifications

### Add support for new AST node type
1. Edit `parser.go`
2. Add case in appropriate function (`declToNode`, `stmtToNode`, `exprToNode`)
3. Return `*ASTNode` with appropriate `Label`, `IndentLevel`, `Children`

### Modify tree display
1. Edit `right_panel.go`
2. Modify `buildListItems()` for display format
3. Modify `toggleNodeCollapse()` for interaction

### Change text input behavior
1. Edit `left_panel.go`
2. Use `TextInput` methods: `SetMultiline()`, `SetAutoWrap()`, `SetTabular()`

## Dependencies

Do not add unnecessary dependencies. Current deps:
- `github.com/guigui-gui/guigui` - GUI framework
- `github.com/hajimehoshi/ebiten/v2` - Graphics (transitive)
- `golang.org/x/tools/txtar` - txtar parsing
- `golang.org/x/text/language` - Locale support (transitive)

## Error Handling

- Parse errors should be displayed in the right panel, not crash the app
- Use `parseErr` field in `RightPanel` to store and display errors

## Performance Considerations

- `Build()` is called frequently; avoid expensive operations
- Use `initialized` flag pattern for one-time setup
- AST parsing happens on button click, not on every keystroke
