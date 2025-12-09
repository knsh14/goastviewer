// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

const defaultSource = `-- main.go --
package main

import "fmt"

// Person represents a person with name and age.
type Person struct {
	Name string
	Age  int
}

// Greet returns a greeting message.
func (p *Person) Greet() string {
	return fmt.Sprintf("Hello, I'm %s", p.Name)
}

func main() {
	p := &Person{
		Name: "Alice",
		Age:  30,
	}
	fmt.Println(p.Greet())
}
`

type LeftPanel struct {
	guigui.DefaultWidget

	titleText   basicwidget.Text
	textInput   basicwidget.TextInput
	parseButton basicwidget.Button

	onSourceChanged func(string)
	currentSource   string
	initialized     bool
}

func (l *LeftPanel) SetOnSourceChanged(f func(string)) {
	l.onSourceChanged = f
}

func (l *LeftPanel) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&l.titleText)
	adder.AddChild(&l.textInput)
	adder.AddChild(&l.parseButton)

	l.titleText.SetValue("txtar Format Go Code:")
	l.titleText.SetBold(true)

	l.textInput.SetMultiline(true)
	l.textInput.SetAutoWrap(false)
	l.textInput.SetTabular(true)
	l.textInput.SetVerticalAlign(basicwidget.VerticalAlignTop)
	l.textInput.SetHorizontalAlign(basicwidget.HorizontalAlignStart)

	// Initialize with default source only once
	if !l.initialized {
		l.initialized = true
		l.currentSource = defaultSource
		l.textInput.SetValue(defaultSource)
		if l.onSourceChanged != nil {
			l.onSourceChanged(defaultSource)
		}
	}

	l.textInput.SetOnValueChanged(func(text string, committed bool) {
		l.currentSource = text
	})

	l.parseButton.SetText("Parse AST")
	l.parseButton.SetOnDown(func() {
		if l.onSourceChanged != nil {
			l.onSourceChanged(l.currentSource)
		}
	})

	return nil
}

func (l *LeftPanel) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	bounds := widgetBounds.Bounds()

	titleSize := l.titleText.Measure(context, guigui.FixedWidthConstraints(bounds.Dx()-u))
	buttonSize := l.parseButton.Measure(context, guigui.Constraints{})

	// Title at top
	titleBounds := image.Rectangle{
		Min: bounds.Min.Add(image.Pt(u/2, u/2)),
		Max: bounds.Min.Add(image.Pt(bounds.Dx()-u/2, u/2+titleSize.Y)),
	}
	layouter.LayoutWidget(&l.titleText, titleBounds)

	// Button at bottom
	buttonBounds := image.Rectangle{
		Min: image.Pt(bounds.Min.X+u/2, bounds.Max.Y-u/2-buttonSize.Y),
		Max: image.Pt(bounds.Max.X-u/2, bounds.Max.Y-u/2),
	}
	layouter.LayoutWidget(&l.parseButton, buttonBounds)

	// TextInput in the middle
	textBounds := image.Rectangle{
		Min: image.Pt(bounds.Min.X+u/2, titleBounds.Max.Y+u/2),
		Max: image.Pt(bounds.Max.X-u/2, buttonBounds.Min.Y-u/2),
	}
	layouter.LayoutWidget(&l.textInput, textBounds)
}
