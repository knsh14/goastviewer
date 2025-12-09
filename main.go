// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"image"
	"os"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/guigui-gui/guigui/basicwidget/cjkfont"
)

type Root struct {
	guigui.DefaultWidget

	background    basicwidget.Background
	leftPanel     LeftPanel
	rightPanel    RightPanel

	locales           []language.Tag
	faceSourceEntries []basicwidget.FaceSourceEntry
}

func (r *Root) updateFontFaceSources(context *guigui.Context) {
	r.locales = slices.Delete(r.locales, 0, len(r.locales))
	r.locales = context.AppendLocales(r.locales)
	r.faceSourceEntries = slices.Delete(r.faceSourceEntries, 0, len(r.faceSourceEntries))
	r.faceSourceEntries = cjkfont.AppendRecommendedFaceSourceEntries(r.faceSourceEntries, r.locales)
	basicwidget.SetFaceSources(r.faceSourceEntries)
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&r.background)
	adder.AddChild(&r.leftPanel)
	adder.AddChild(&r.rightPanel)

	r.updateFontFaceSources(context)
	r.leftPanel.SetOnSourceChanged(func(source string) {
		r.rightPanel.SetSource(source)
	})
	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&r.background, widgetBounds.Bounds())
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &r.leftPanel,
				Size:   guigui.FlexibleSize(1),
			},
			{
				Widget: &r.rightPanel,
				Size:   guigui.FlexibleSize(1),
			},
		},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func main() {
	op := &guigui.RunOptions{
		Title:      "Go AST Viewer",
		WindowSize: image.Pt(1200, 800),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
