// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image/color"

	"github.com/xiaoshengduan/gio-fly/internal/f32color"
	"github.com/xiaoshengduan/gio-fly/io/semantic"
	"github.com/xiaoshengduan/gio-fly/layout"
	"github.com/xiaoshengduan/gio-fly/op"
	"github.com/xiaoshengduan/gio-fly/op/paint"
	"github.com/xiaoshengduan/gio-fly/text"
	"github.com/xiaoshengduan/gio-fly/unit"
	"github.com/xiaoshengduan/gio-fly/widget"
)

type EditorStyle struct {
	Font     text.Font
	TextSize unit.Sp
	// Color is the text color.
	Color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.NRGBA
	// SelectionColor is the color of the background for selected text.
	SelectionColor color.NRGBA
	Editor         *widget.Editor

	shaper text.Shaper
}

func Editor(th *Theme, editor *widget.Editor, hint string) EditorStyle {
	return EditorStyle{
		Editor:         editor,
		TextSize:       th.TextSize,
		Color:          th.Palette.Fg,
		shaper:         th.Shaper,
		Hint:           hint,
		HintColor:      f32color.MulAlpha(th.Palette.Fg, 0xbb),
		SelectionColor: f32color.MulAlpha(th.Palette.ContrastBg, 0x60),
	}
}

func (e EditorStyle) Layout(gtx layout.Context) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	var maxlines int
	if e.Editor.SingleLine {
		maxlines = 1
	}
	tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: maxlines}
	dims := tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = e.Editor.Layout(gtx, e.shaper, e.Font, e.TextSize, func(gtx layout.Context) layout.Dimensions {
		semantic.Editor.Add(gtx.Ops)
		disabled := gtx.Queue == nil
		if e.Editor.Len() > 0 {
			paint.ColorOp{Color: blendDisabledColor(disabled, e.SelectionColor)}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: blendDisabledColor(disabled, e.Color)}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled {
			paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	})
	return dims
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return f32color.Disabled(c)
	}
	return c
}
