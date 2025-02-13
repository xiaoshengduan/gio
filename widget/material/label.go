// SPDX-License-Identifier: Unlicense OR MIT

package material

import (
	"image/color"

	"github.com/xiaoshengduan/gio-fly/layout"
	"github.com/xiaoshengduan/gio-fly/op/paint"
	"github.com/xiaoshengduan/gio-fly/text"
	"github.com/xiaoshengduan/gio-fly/unit"
	"github.com/xiaoshengduan/gio-fly/widget"
)

type LabelStyle struct {
	// Face defines the text style.
	Font text.Font
	// Color is the text color.
	Color color.NRGBA
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Sp

	shaper text.Shaper
}

func H1(th *Theme, txt string) LabelStyle {
	label := Label(th, th.TextSize*96.0/16.0, txt)
	label.Font.Weight = text.Light
	return label
}

func H2(th *Theme, txt string) LabelStyle {
	label := Label(th, th.TextSize*60.0/16.0, txt)
	label.Font.Weight = text.Light
	return label
}

func H3(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*48.0/16.0, txt)
}

func H4(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*34.0/16.0, txt)
}

func H5(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*24.0/16.0, txt)
}

func H6(th *Theme, txt string) LabelStyle {
	label := Label(th, th.TextSize*20.0/16.0, txt)
	label.Font.Weight = text.Medium
	return label
}

func Subtitle1(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*16.0/16.0, txt)
}

func Subtitle2(th *Theme, txt string) LabelStyle {
	label := Label(th, th.TextSize*14.0/16.0, txt)
	label.Font.Weight = text.Medium
	return label
}

func Body1(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize, txt)
}

func Body2(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*14.0/16.0, txt)
}

func Caption(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*12.0/16.0, txt)
}

func Overline(th *Theme, txt string) LabelStyle {
	return Label(th, th.TextSize*10.0/16.0, txt)
}

func Label(th *Theme, size unit.Sp, txt string) LabelStyle {
	return LabelStyle{
		Text:     txt,
		Color:    th.Palette.Fg,
		TextSize: size,
		shaper:   th.Shaper,
	}
}

func (l LabelStyle) Layout(gtx layout.Context) layout.Dimensions {
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	tl := widget.Label{Alignment: l.Alignment, MaxLines: l.MaxLines}
	return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
