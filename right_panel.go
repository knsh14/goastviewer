// SPDX-License-Identifier: Apache-2.0

package main

import (
	"image"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

type RightPanel struct {
	guigui.DefaultWidget

	panel     basicwidget.Panel
	titleText basicwidget.Text
	treeList  basicwidget.List[int]
	errorText basicwidget.Text

	source    string
	astNodes  []*ASTNode
	listItems []basicwidget.ListItem[int]
	parseErr  error
}

func (r *RightPanel) SetSource(source string) {
	r.source = source
	r.parseAST()
}

func (r *RightPanel) parseAST() {
	if r.source == "" {
		r.astNodes = nil
		r.parseErr = nil
		return
	}

	nodes, err := ParseTxtar(r.source)
	if err != nil {
		r.parseErr = err
		r.astNodes = nil
		return
	}

	r.parseErr = nil
	r.astNodes = nodes
}

func (r *RightPanel) buildListItems() {
	r.listItems = r.listItems[:0]

	if r.astNodes == nil {
		return
	}

	flatNodes := FlattenNodes(r.astNodes)
	for i, node := range flatNodes {
		hasChildren := len(node.Children) > 0
		label := node.Label
		if hasChildren {
			if node.Collapsed {
				label = "[+] " + label
			} else {
				label = "[-] " + label
			}
		} else {
			label = "    " + label
		}

		r.listItems = append(r.listItems, basicwidget.ListItem[int]{
			Text:        label,
			IndentLevel: node.IndentLevel,
			Value:       i,
			Collapsed:   node.Collapsed,
		})
	}
}

func (r *RightPanel) toggleNodeCollapse(index int) {
	if index < 0 {
		return
	}

	flatNodes := FlattenNodes(r.astNodes)
	if index >= len(flatNodes) {
		return
	}

	node := flatNodes[index]
	if len(node.Children) > 0 {
		node.Collapsed = !node.Collapsed
	}
}

func (r *RightPanel) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&r.panel)
	r.panel.SetContent(&rightPanelContent{rightPanel: r})
	r.panel.SetAutoBorder(true)
	r.panel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)

	return nil
}

func (r *RightPanel) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	layouter.LayoutWidget(&r.panel, widgetBounds.Bounds())
}

type rightPanelContent struct {
	guigui.DefaultWidget
	rightPanel *RightPanel
}

func (p *rightPanelContent) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddChild(&p.rightPanel.titleText)

	p.rightPanel.titleText.SetValue("AST Tree:")
	p.rightPanel.titleText.SetBold(true)

	if p.rightPanel.parseErr != nil {
		adder.AddChild(&p.rightPanel.errorText)
		p.rightPanel.errorText.SetValue("Error: " + p.rightPanel.parseErr.Error())
	} else {
		adder.AddChild(&p.rightPanel.treeList)
		p.rightPanel.buildListItems()
		p.rightPanel.treeList.SetItems(p.rightPanel.listItems)
		p.rightPanel.treeList.SetStripeVisible(true)
		p.rightPanel.treeList.SetOnItemSelected(func(index int) {
			p.rightPanel.toggleNodeCollapse(index)
		})
		p.rightPanel.treeList.SetOnItemExpanderToggled(func(index int, expanded bool) {
			flatNodes := FlattenNodes(p.rightPanel.astNodes)
			if index < len(flatNodes) {
				flatNodes[index].Collapsed = !expanded
			}
		})
	}

	return nil
}

func (p *rightPanelContent) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)
	bounds := widgetBounds.Bounds()

	var contentWidget guigui.Widget
	if p.rightPanel.parseErr != nil {
		contentWidget = &p.rightPanel.errorText
	} else {
		contentWidget = &p.rightPanel.treeList
	}

	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items: []guigui.LinearLayoutItem{
			{
				Widget: &p.rightPanel.titleText,
			},
			{
				Widget: contentWidget,
				Size:   guigui.FlexibleSize(1),
			},
		},
		Gap: u / 2,
		Padding: guigui.Padding{
			Start:  u / 2,
			Top:    u / 2,
			End:    u / 2,
			Bottom: u / 2,
		},
	}).LayoutWidgets(context, bounds, layouter)
}

func (p *rightPanelContent) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	u := basicwidget.UnitSize(context)
	if w, ok := constraints.FixedWidth(); ok {
		return image.Pt(w, 20*u)
	}
	return image.Pt(20*u, 20*u)
}
