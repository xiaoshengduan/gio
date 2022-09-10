// SPDX-License-Identifier: Unlicense OR MIT

package gesture

import (
	"image"
	"testing"
	"time"

	"github.com/xiaoshengduan/gio-fly/f32"
	"github.com/xiaoshengduan/gio-fly/io/event"
	"github.com/xiaoshengduan/gio-fly/io/pointer"
	"github.com/xiaoshengduan/gio-fly/io/router"
	"github.com/xiaoshengduan/gio-fly/op"
	"github.com/xiaoshengduan/gio-fly/op/clip"
)

func TestHover(t *testing.T) {
	ops := new(op.Ops)
	var h Hover
	rect := image.Rect(20, 20, 40, 40)
	stack := clip.Rect(rect).Push(ops)
	h.Add(ops)
	stack.Pop()
	r := new(router.Router)
	r.Frame(ops)

	r.Queue(
		pointer.Event{Type: pointer.Move, Position: f32.Pt(30, 30)},
	)
	if !h.Hovered(r) {
		t.Fatal("expected hovered")
	}

	r.Queue(
		pointer.Event{Type: pointer.Move, Position: f32.Pt(50, 50)},
	)
	if h.Hovered(r) {
		t.Fatal("expected not hovered")
	}
}

func TestMouseClicks(t *testing.T) {
	for _, tc := range []struct {
		label  string
		events []event.Event
		clicks []int // number of combined clicks per click (single, double...)
	}{
		{
			label:  "single click",
			events: mouseClickEvents(200 * time.Millisecond),
			clicks: []int{1},
		},
		{
			label: "double click",
			events: mouseClickEvents(
				100*time.Millisecond,
				100*time.Millisecond+doubleClickDuration-1),
			clicks: []int{1, 2},
		},
		{
			label: "two single clicks",
			events: mouseClickEvents(
				100*time.Millisecond,
				100*time.Millisecond+doubleClickDuration+1),
			clicks: []int{1, 1},
		},
	} {
		t.Run(tc.label, func(t *testing.T) {
			var click Click
			var ops op.Ops
			click.Add(&ops)

			var r router.Router
			r.Frame(&ops)
			r.Queue(tc.events...)

			events := click.Events(&r)
			clicks := filterMouseClicks(events)
			if got, want := len(clicks), len(tc.clicks); got != want {
				t.Fatalf("got %d mouse clicks, expected %d", got, want)
			}

			for i, click := range clicks {
				if got, want := click.NumClicks, tc.clicks[i]; got != want {
					t.Errorf("got %d combined mouse clicks, expected %d", got, want)
				}
			}
		})
	}
}

func mouseClickEvents(times ...time.Duration) []event.Event {
	press := pointer.Event{
		Type:    pointer.Press,
		Source:  pointer.Mouse,
		Buttons: pointer.ButtonPrimary,
	}
	events := make([]event.Event, 0, 2*len(times))
	for _, t := range times {
		release := press
		release.Type = pointer.Release
		release.Time = t
		events = append(events, press, release)
	}
	return events
}

func filterMouseClicks(events []ClickEvent) []ClickEvent {
	var clicks []ClickEvent
	for _, ev := range events {
		if ev.Type == TypeClick {
			clicks = append(clicks, ev)
		}
	}
	return clicks
}
