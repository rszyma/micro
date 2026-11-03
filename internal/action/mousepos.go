package action

import (
	"github.com/gdamore/tcell/v3"
	"github.com/micro-editor/micro/v2/internal/display"
	"github.com/micro-editor/micro/v2/internal/screen"
)

// mouseEventCells maps a tcell mouse event from pixels to cells when pixel
// mouse mode is active. It converts absolute screen pixels to absolute
// screen cells using the window origin (no gutter, no clamping): the
// caller owns the geometry. fracX holds the fractional position
// inside the clicked cell (0 = left edge, 1 = right edge).
// In cell mode the fraction is 0 and the position passes through.
func mouseEventCells(e *tcell.EventMouse, winX, winY int) (mx, my int, fracX float64) {
	mx, my = e.Position()
	if !screen.MousePixel {
		return mx, my, 0
	}
	cw, ch := screen.Screen.MousePixelSize()
	if cw <= 0 || ch <= 0 {
		return mx, my, 0
	}
	relX := mx - winX*cw
	relY := my - winY*ch
	if relX < 0 {
		relX = 0
	}
	if relY < 0 {
		relY = 0
	}
	cellX := relX / cw
	cellY := relY / ch
	fracX = float64(relX%cw) / float64(cw)
	return winX + cellX, winY + cellY, fracX
}

// mouseCellPos maps a tcell mouse event to a cell position plus the
// fractional position inside the clicked cell (0 = left edge,
// 1 = right edge). In pixel mode pixel offsets are divided by the
// cell size; in cell mode the fraction is 0 and the position passes
// through unchanged. The window origin and horizontal gutter come from
// ViewMetrics so clicks in the gutter clamp to the nearest buffer cell
// instead of producing a negative offset.
func mouseCellPos(e *tcell.EventMouse, w display.BWindow) (mx, my int, fracX float64) {
	mx, my = e.Position()
	if !screen.MousePixel {
		return mx, my, 0
	}
	cw, ch := screen.Screen.MousePixelSize()
	if cw <= 0 || ch <= 0 {
		return mx, my, 0
	}
	// Mouse positions are in screen coordinates. The buffer view starts
	// at (winX+gutter, winY); pixels outside it (gutter, dividers, other
	// panes) clamp to the nearest edge of the view.
	winX, gutter, _ := w.ViewMetrics()
	bv := w.BufView()
	relX := mx - (winX+gutter)*cw
	relY := my - bv.Y*ch
	if relX < 0 {
		relX = 0
	}
	if relY < 0 {
		relY = 0
	}
	cellX := relX / cw
	cellY := relY / ch
	fracX = float64(relX%cw) / float64(cw)
	mx = winX + gutter + cellX
	my = bv.Y + cellY
	return mx, my, fracX
}
