// SPDX-License-Identifier: Unlicense OR MIT

package rendertest

import (
	"image"
	"image/color"
	"math"
	"testing"

	"golang.org/x/image/colornames"

	"gioui.org/f32"
	"gioui.org/internal/f32color"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func TestTransformMacro(t *testing.T) {
	// testcase resulting from original bug when rendering layout.Stacked

	// Build clip-path.
	c := constSqPath()

	run(t, func(o *op.Ops) {

		// render the first Stacked item
		m1 := op.Record(o)
		dr := image.Rect(0, 0, 128, 50)
		paint.FillShape(o, black, clip.Rect(dr).Op())
		c1 := m1.Stop()

		// Render the second stacked item
		m2 := op.Record(o)
		paint.ColorOp{Color: red}.Add(o)
		// Simulate a draw text call
		t := op.Offset(f32.Pt(0, 10)).Push(o)

		// Apply the clip-path.
		cl := c.Push(o)

		paint.PaintOp{}.Add(o)
		cl.Pop()
		t.Pop()

		c2 := m2.Stop()

		// Call each of them in a transform
		t = op.Offset(f32.Pt(0, 0)).Push(o)
		c1.Add(o)
		t.Pop()
		t = op.Offset(f32.Pt(0, 0)).Push(o)
		c2.Add(o)
		t.Pop()
	}, func(r result) {
		r.expect(5, 15, colornames.Red)
		r.expect(15, 15, colornames.Black)
		r.expect(11, 51, transparent)
	})
}

func TestRepeatedPaintsZ(t *testing.T) {
	run(t, func(o *op.Ops) {
		// Draw a rectangle
		paint.FillShape(o, black, clip.Rect(image.Rect(0, 0, 128, 50)).Op())

		builder := clip.Path{}
		builder.Begin(o)
		builder.Move(f32.Pt(0, 0))
		builder.Line(f32.Pt(10, 0))
		builder.Line(f32.Pt(0, 10))
		builder.Line(f32.Pt(-10, 0))
		builder.Line(f32.Pt(0, -10))
		p := builder.End()
		defer clip.Outline{
			Path: p,
		}.Op().Push(o).Pop()
		paint.Fill(o, red)
	}, func(r result) {
		r.expect(5, 5, colornames.Red)
		r.expect(11, 15, colornames.Black)
		r.expect(11, 51, transparent)
	})
}

func TestNoClipFromPaint(t *testing.T) {
	// ensure that a paint operation does not pollute the state
	// by leaving any clip paths in place.
	run(t, func(o *op.Ops) {
		a := f32.Affine2D{}.Rotate(f32.Pt(20, 20), math.Pi/4)
		defer op.Affine(a).Push(o).Pop()
		paint.FillShape(o, red, clip.Rect(image.Rect(10, 10, 30, 30)).Op())
		a = f32.Affine2D{}.Rotate(f32.Pt(20, 20), -math.Pi/4)
		defer op.Affine(a).Push(o).Pop()

		paint.FillShape(o, black, clip.Rect(image.Rect(0, 0, 50, 50)).Op())
	}, func(r result) {
		r.expect(1, 1, colornames.Black)
		r.expect(20, 20, colornames.Black)
		r.expect(49, 49, colornames.Black)
		r.expect(51, 51, transparent)
	})
}

func TestDeferredPaint(t *testing.T) {
	run(t, func(o *op.Ops) {
		cl := clip.Rect(image.Rect(0, 0, 80, 80)).Op().Push(o)
		paint.ColorOp{Color: color.NRGBA{A: 0x60, G: 0xff}}.Add(o)
		paint.PaintOp{}.Add(o)
		cl.Pop()

		t := op.Affine(f32.Affine2D{}.Offset(f32.Pt(20, 20))).Push(o)
		m := op.Record(o)
		cl2 := clip.Rect(image.Rect(0, 0, 80, 80)).Op().Push(o)
		paint.ColorOp{Color: color.NRGBA{A: 0x60, R: 0xff, G: 0xff}}.Add(o)
		paint.PaintOp{}.Add(o)
		cl2.Pop()
		paintMacro := m.Stop()
		op.Defer(o, paintMacro)
		t.Pop()

		defer op.Affine(f32.Affine2D{}.Offset(f32.Pt(10, 10))).Push(o).Pop()
		defer clip.Rect(image.Rect(0, 0, 80, 80)).Op().Push(o).Pop()
		paint.ColorOp{Color: color.NRGBA{A: 0x60, B: 0xff}}.Add(o)
		paint.PaintOp{}.Add(o)
	}, func(r result) {
	})
}

func constSqPath() clip.Op {
	innerOps := new(op.Ops)
	builder := clip.Path{}
	builder.Begin(innerOps)
	builder.Move(f32.Pt(0, 0))
	builder.Line(f32.Pt(10, 0))
	builder.Line(f32.Pt(0, 10))
	builder.Line(f32.Pt(-10, 0))
	builder.Line(f32.Pt(0, -10))
	p := builder.End()
	return clip.Outline{Path: p}.Op()
}

func constSqCirc() clip.Op {
	innerOps := new(op.Ops)
	return clip.RRect{Rect: f32.Rect(0, 0, 40, 40),
		NW: 20, NE: 20, SW: 20, SE: 20}.Op(innerOps)
}

func drawChild(ops *op.Ops, text clip.Op) op.CallOp {
	r1 := op.Record(ops)
	cl := text.Push(ops)
	paint.PaintOp{}.Add(ops)
	cl.Pop()
	return r1.Stop()
}

func TestReuseStencil(t *testing.T) {
	txt := constSqPath()
	run(t, func(ops *op.Ops) {
		c1 := drawChild(ops, txt)
		c2 := drawChild(ops, txt)

		// lay out the children
		c1.Add(ops)

		defer op.Offset(f32.Pt(0, 50)).Push(ops).Pop()
		c2.Add(ops)
	}, func(r result) {
		r.expect(5, 5, colornames.Black)
		r.expect(5, 55, colornames.Black)
	})
}

func TestBuildOffscreen(t *testing.T) {
	// Check that something we in one frame build outside the screen
	// still is rendered correctly if moved into the screen in a later
	// frame.

	txt := constSqCirc()
	draw := func(off float32, o *op.Ops) {
		defer op.Offset(f32.Pt(0, off)).Push(o).Pop()
		defer txt.Push(o).Pop()
		paint.PaintOp{}.Add(o)
	}

	multiRun(t,
		frame(
			func(ops *op.Ops) {
				draw(-100, ops)
			}, func(r result) {
				r.expect(5, 5, transparent)
				r.expect(20, 20, transparent)
			}),
		frame(
			func(ops *op.Ops) {
				draw(0, ops)
			}, func(r result) {
				r.expect(2, 2, transparent)
				r.expect(20, 20, colornames.Black)
				r.expect(38, 38, transparent)
			}))
}

func TestNegativeOverlaps(t *testing.T) {
	run(t, func(ops *op.Ops) {
		defer clip.RRect{Rect: f32.Rect(50, 50, 100, 100)}.Push(ops).Pop()
		clip.Rect(image.Rect(0, 120, 100, 122)).Push(ops).Pop()
		paint.PaintOp{}.Add(ops)
	}, func(r result) {
		r.expect(60, 60, transparent)
		r.expect(60, 110, transparent)
		r.expect(60, 120, transparent)
		r.expect(60, 122, transparent)
	})
}

func TestDepthOverlap(t *testing.T) {
	run(t, func(ops *op.Ops) {
		paint.FillShape(ops, red, clip.Rect{Max: image.Pt(128, 64)}.Op())
		paint.FillShape(ops, green, clip.Rect{Max: image.Pt(64, 128)}.Op())
	}, func(r result) {
		r.expect(96, 32, colornames.Red)
		r.expect(32, 96, colornames.Green)
		r.expect(32, 32, colornames.Green)
	})
}

type Gradient struct {
	From, To color.NRGBA
}

var gradients = []Gradient{
	{From: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, To: color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}},
	{From: color.NRGBA{R: 0x19, G: 0xFF, B: 0x19, A: 0xFF}, To: color.NRGBA{R: 0xFF, G: 0x19, B: 0x19, A: 0xFF}},
	{From: color.NRGBA{R: 0xFF, G: 0x19, B: 0x19, A: 0xFF}, To: color.NRGBA{R: 0x19, G: 0x19, B: 0xFF, A: 0xFF}},
	{From: color.NRGBA{R: 0x19, G: 0x19, B: 0xFF, A: 0xFF}, To: color.NRGBA{R: 0x19, G: 0xFF, B: 0x19, A: 0xFF}},
	{From: color.NRGBA{R: 0x19, G: 0xFF, B: 0xFF, A: 0xFF}, To: color.NRGBA{R: 0xFF, G: 0x19, B: 0x19, A: 0xFF}},
	{From: color.NRGBA{R: 0xFF, G: 0xFF, B: 0x19, A: 0xFF}, To: color.NRGBA{R: 0x19, G: 0x19, B: 0xFF, A: 0xFF}},
}

func TestLinearGradient(t *testing.T) {
	t.Skip("linear gradients don't support transformations")

	const gradienth = 8
	// 0.5 offset from ends to ensure that the center of the pixel
	// aligns with gradient from and to colors.
	pixelAligned := f32.Rect(0.5, 0, 127.5, gradienth)
	samples := []int{0, 12, 32, 64, 96, 115, 127}

	run(t, func(ops *op.Ops) {
		gr := f32.Rect(0, 0, 128, gradienth)
		for _, g := range gradients {
			paint.LinearGradientOp{
				Stop1:  f32.Pt(gr.Min.X, gr.Min.Y),
				Color1: g.From,
				Stop2:  f32.Pt(gr.Max.X, gr.Min.Y),
				Color2: g.To,
			}.Add(ops)
			cl := clip.RRect{Rect: gr}.Push(ops)
			t1 := op.Affine(f32.Affine2D{}.Offset(pixelAligned.Min)).Push(ops)
			t2 := scale(pixelAligned.Dx()/128, 1).Push(ops)
			paint.PaintOp{}.Add(ops)
			t2.Pop()
			t1.Pop()
			cl.Pop()
			gr = gr.Add(f32.Pt(0, gradienth))
		}
	}, func(r result) {
		gr := pixelAligned
		for _, g := range gradients {
			from := f32color.LinearFromSRGB(g.From)
			to := f32color.LinearFromSRGB(g.To)
			for _, p := range samples {
				exp := lerp(from, to, float32(p)/float32(r.img.Bounds().Dx()-1))
				r.expect(p, int(gr.Min.Y+gradienth/2), f32color.NRGBAToRGBA(exp.SRGB()))
			}
			gr = gr.Add(f32.Pt(0, gradienth))
		}
	})
}

func TestLinearGradientAngled(t *testing.T) {
	run(t, func(ops *op.Ops) {
		paint.LinearGradientOp{
			Stop1:  f32.Pt(64, 64),
			Color1: black,
			Stop2:  f32.Pt(0, 0),
			Color2: red,
		}.Add(ops)
		cl := clip.Rect(image.Rect(0, 0, 64, 64)).Push(ops)
		paint.PaintOp{}.Add(ops)
		cl.Pop()

		paint.LinearGradientOp{
			Stop1:  f32.Pt(64, 64),
			Color1: white,
			Stop2:  f32.Pt(128, 0),
			Color2: green,
		}.Add(ops)
		cl = clip.Rect(image.Rect(64, 0, 128, 64)).Push(ops)
		paint.PaintOp{}.Add(ops)
		cl.Pop()

		paint.LinearGradientOp{
			Stop1:  f32.Pt(64, 64),
			Color1: black,
			Stop2:  f32.Pt(128, 128),
			Color2: blue,
		}.Add(ops)
		cl = clip.Rect(image.Rect(64, 64, 128, 128)).Push(ops)
		paint.PaintOp{}.Add(ops)
		cl.Pop()

		paint.LinearGradientOp{
			Stop1:  f32.Pt(64, 64),
			Color1: white,
			Stop2:  f32.Pt(0, 128),
			Color2: magenta,
		}.Add(ops)
		cl = clip.Rect(image.Rect(0, 64, 64, 128)).Push(ops)
		paint.PaintOp{}.Add(ops)
		cl.Pop()
	}, func(r result) {})
}

func TestZeroImage(t *testing.T) {
	ops := new(op.Ops)
	w := newWindow(t, 10, 10)
	paint.ImageOp{}.Add(ops)
	paint.PaintOp{}.Add(ops)
	if err := w.Frame(ops); err != nil {
		t.Error(err)
	}
}

func TestImageRGBA(t *testing.T) {
	run(t, func(o *op.Ops) {
		w := newWindow(t, 10, 10)

		im := image.NewRGBA(image.Rect(0, 0, 5, 5))
		im.Set(3, 3, colornames.Black)
		im.Set(4, 3, colornames.Black)
		im.Set(3, 4, colornames.Black)
		im.Set(4, 4, colornames.Black)
		im = im.SubImage(image.Rect(2, 2, 5, 5)).(*image.RGBA)
		paint.NewImageOp(im).Add(o)
		paint.PaintOp{}.Add(o)
		if err := w.Frame(o); err != nil {
			t.Error(err)
		}
	}, func(r result) {
		r.expect(1, 1, colornames.Black)
		r.expect(2, 1, colornames.Black)
		r.expect(1, 2, colornames.Black)
		r.expect(2, 2, colornames.Black)
	})
}

// lerp calculates linear interpolation with color b and p.
func lerp(a, b f32color.RGBA, p float32) f32color.RGBA {
	return f32color.RGBA{
		R: a.R*(1-p) + b.R*p,
		G: a.G*(1-p) + b.G*p,
		B: a.B*(1-p) + b.B*p,
		A: a.A*(1-p) + b.A*p,
	}
}
