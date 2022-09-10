// SPDX-License-Identifier: Unlicense OR MIT

package widget_test

import (
	"image"
	"testing"

	"github.com/xiaoshengduan/gio-fly/f32"
	"github.com/xiaoshengduan/gio-fly/io/pointer"
	"github.com/xiaoshengduan/gio-fly/io/router"
	"github.com/xiaoshengduan/gio-fly/io/semantic"
	"github.com/xiaoshengduan/gio-fly/io/system"
	"github.com/xiaoshengduan/gio-fly/layout"
	"github.com/xiaoshengduan/gio-fly/op"
	"github.com/xiaoshengduan/gio-fly/widget"
)

func TestBool(t *testing.T) {
	var (
		ops op.Ops
		r   router.Router
		b   widget.Bool
	)
	gtx := layout.NewContext(&ops, system.FrameEvent{Queue: &r})
	layout := func() {
		b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.CheckBox.Add(gtx.Ops)
			semantic.DescriptionOp("description").Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(100, 100)}
		})
	}
	layout()
	r.Frame(gtx.Ops)
	r.Queue(
		pointer.Event{
			Source:   pointer.Touch,
			Type:     pointer.Press,
			Position: f32.Pt(50, 50),
		},
		pointer.Event{
			Source:   pointer.Touch,
			Type:     pointer.Release,
			Position: f32.Pt(50, 50),
		},
	)
	ops.Reset()
	layout()
	r.Frame(gtx.Ops)
	tree := r.AppendSemantics(nil)
	n := tree[0].Children[0].Desc
	if n.Description != "description" {
		t.Errorf("unexpected semantic description: %s", n.Description)
	}
	if n.Class != semantic.CheckBox {
		t.Errorf("unexpected semantic class: %v", n.Class)
	}
	if !b.Value || !n.Selected {
		t.Error("click did not select")
	}
}
