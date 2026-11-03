package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringWidth(t *testing.T) {
	bytes := []byte("\tPot să \tmănânc sticlă și ea nu mă rănește.")

	n := StringWidth(bytes, 23, 4)
	assert.Equal(t, 26, n)
}

func TestGetCharPosInLineFrac(t *testing.T) {
	// single-cell runes: left half stays, right half advances
	line := []byte("ab")
	assert.Equal(t, 0, GetCharPosInLineFrac(line, 0.2, 4))
	assert.Equal(t, 1, GetCharPosInLineFrac(line, 0.7, 4))
	assert.Equal(t, 1, GetCharPosInLineFrac(line, 1.0, 4))
	assert.Equal(t, 1, GetCharPosInLineFrac(line, 1.2, 4))
	assert.Equal(t, 2, GetCharPosInLineFrac(line, 1.7, 4))

	// wide rune (2 cells): left half stays before it, right half past it
	wide := []byte("a\u65e5b") // "a"日"b", 日 is width 2
	assert.Equal(t, 1, GetCharPosInLineFrac(wide, 1.2, 4))
	assert.Equal(t, 2, GetCharPosInLineFrac(wide, 2.7, 4))

	// tab (4 cells): left half stays, right half advances
	tab := []byte("\ta")
	assert.Equal(t, 0, GetCharPosInLineFrac(tab, 1.5, 4))
	assert.Equal(t, 1, GetCharPosInLineFrac(tab, 2.5, 4))

	// integer positions behave exactly like GetCharPosInLine
	for _, tc := range []struct {
		line []byte
		pos  int
	}{
		{line, 1}, {line, 2}, {wide, 3}, {tab, 4},
	} {
		assert.Equal(t, GetCharPosInLine(tc.line, tc.pos, 4), GetCharPosInLineFrac(tc.line, float64(tc.pos), 4))
	}
}

func TestSliceVisualEnd(t *testing.T) {
	s := []byte("\thello")
	slc, n, _ := SliceVisualEnd(s, 2, 4)
	assert.Equal(t, []byte("\thello"), slc)
	assert.Equal(t, 2, n)

	slc, n, _ = SliceVisualEnd(s, 1, 4)
	assert.Equal(t, []byte("\thello"), slc)
	assert.Equal(t, 1, n)

	slc, n, _ = SliceVisualEnd(s, 4, 4)
	assert.Equal(t, []byte("hello"), slc)
	assert.Equal(t, 0, n)

	slc, n, _ = SliceVisualEnd(s, 5, 4)
	assert.Equal(t, []byte("ello"), slc)
	assert.Equal(t, 0, n)
}
