package display

import (
	"testing"

	"github.com/micro-editor/micro/v2/internal/buffer"
	"github.com/micro-editor/micro/v2/internal/config"
	ulua "github.com/micro-editor/micro/v2/internal/lua"
	lua "github.com/yuin/gopher-lua"
)

func init() {
	ulua.L = lua.NewState()
	config.InitRuntimeFiles(false)
	config.InitGlobalSettings()
}

func newFracTestWindow(text string, bufWidth int, wordwrap bool) *BufWindow {
	b := buffer.NewBufferFromString(text, "", buffer.BTDefault)
	b.Settings["softwrap"] = true
	b.Settings["wordwrap"] = wordwrap
	b.Settings["tabsize"] = float64(4)
	b.Settings["statusline"] = true
	b.Settings["scrollbar"] = false
	w := NewBufWindow(0, 0, 80, 24, b)
	w.Resize(80, 24)
	// Decouple from gutter math: LocFromVLoc only uses bufWidth.
	w.bufWidth = bufWidth
	return w
}

func TestLocFromVLocFracWideRune(t *testing.T) {
	// "a中b": a=[0,1) 中=[1,3) b=[3,4), single visual row.
	w := newFracTestWindow("a中b", 40, false)

	cases := []struct {
		name string
		x    int
		frac float64
		want buffer.Loc
	}{
		{"left half of wide rune", 1, 0.2, buffer.Loc{1, 0}},
		{"late left half of wide rune", 1, 0.9, buffer.Loc{1, 0}},
		{"right half of wide rune", 2, 0.1, buffer.Loc{2, 0}},
		{"far right half of wide rune", 2, 0.9, buffer.Loc{2, 0}},
		{"right half of plain cell", 0, 0.9, buffer.Loc{1, 0}},
		{"left half of plain cell", 0, 0.2, buffer.Loc{0, 0}},
	}

	for _, c := range cases {
		got := w.LocFromVLocFrac(VLoc{SLoc{0, 0}, c.x}, c.frac)
		if got != c.want {
			t.Errorf("%s: cell %d frac %v: got %v, want %v", c.name, c.x, c.frac, got, c.want)
		}
	}

	// Legacy behavior snaps left on both cells of the wide rune.
	if got := w.LocFromVLoc(VLoc{SLoc{0, 0}, 1}); got != (buffer.Loc{1, 0}) {
		t.Errorf("legacy cell 1: got %v, want {1 0}", got)
	}
	if got := w.LocFromVLoc(VLoc{SLoc{0, 0}, 2}); got != (buffer.Loc{1, 0}) {
		t.Errorf("legacy cell 2: got %v, want {1 0}", got)
	}
}

func TestLocFromVLocFracSoftwrap(t *testing.T) {
	// "ab中de" with bufWidth 4: row0 = a b 中, row1 = d e.
	// 中 spans row0 cells [2,4), d is row1 cell 0.
	w := newFracTestWindow("ab中de", 4, false)

	cases := []struct {
		name string
		row  int
		x    int
		frac float64
		want buffer.Loc
	}{
		{"wrapped rune right cell right half", 1, 0, 0.7, buffer.Loc{4, 0}},
		{"wrapped rune right cell left half", 1, 0, 0.2, buffer.Loc{3, 0}},
		{"wide rune left half on row 0", 0, 2, 0.2, buffer.Loc{2, 0}},
		{"wide rune right half on row 0", 0, 3, 0.1, buffer.Loc{3, 0}},
	}

	for _, c := range cases {
		got := w.LocFromVLocFrac(VLoc{SLoc{0, c.row}, c.x}, c.frac)
		if got != c.want {
			t.Errorf("%s: row %d cell %d frac %v: got %v, want %v", c.name, c.row, c.x, c.frac, got, c.want)
		}
	}

	// A frac of 0 behaves exactly like the legacy path.
	for row := 0; row <= 1; row++ {
		for x := 0; x <= 5; x++ {
			v := VLoc{SLoc{0, row}, x}
			if got, want := w.LocFromVLocFrac(v, 0), w.LocFromVLoc(v); got != want {
				t.Errorf("frac 0 differs from legacy at row %d cell %d: got %v, want %v", row, x, got, want)
			}
		}
	}
}

func TestLocFromVLocFracSoftwrapWordwrap(t *testing.T) {
	// "aa bb cc" with bufWidth 6: row0 = "aa bb ", row1 = "cc".
	w := newFracTestWindow("aa bb cc", 6, true)

	got := w.LocFromVLocFrac(VLoc{SLoc{0, 1}, 0}, 0.9)
	if got != (buffer.Loc{7, 0}) {
		t.Errorf("wrapped word click: got %v, want {7 0}", got)
	}
	if got := w.LocFromVLoc(VLoc{SLoc{0, 1}, 0}); got != (buffer.Loc{6, 0}) {
		t.Errorf("legacy wrapped word click: got %v, want {6 0}", got)
	}

	// Sweep: results stay inside the buffer, no panics.
	n := len([]rune("aa bb cc"))
	for row := 0; row <= 2; row++ {
		for x := 0; x <= 7; x++ {
			for _, frac := range []float64{0, 0.25, 0.5, 0.75, 0.99} {
				got := w.LocFromVLocFrac(VLoc{SLoc{0, row}, x}, frac)
				if got.Y != 0 || got.X < 0 || got.X > n {
					t.Errorf("out of range at row %d cell %d frac %v: got %v", row, x, frac, got)
				}
			}
		}
	}
}
