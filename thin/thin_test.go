package thin

import (
	"fmt"
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ernyoke/imger/imgio"
)

func TestNeighbor(t *testing.T) {
	const (
		allNeighbors = `
@@@
@@@
@@@
`
		noNeighbors = `
---
---
---
`
	)

	img := bitmapToImage(allNeighbors)
	for p := 1; p < 10; p++ {
		n, _ := neighbor(img, 1, 1, p)
		require.Equal(t, Black, n,
			fmt.Sprintf("neighbor %d; expected Black, got White", p))
	}

	// We also check that all neighbors outside the image return as White here.
	img = bitmapToImage(noNeighbors)
	for p := 2; p < 10; p++ {
		for x := 0; x < img.Bounds().Size().X; x++ {
			for y := 0; y < img.Bounds().Size().Y; y++ {
				n, _ := neighbor(img, x, y, p)
				require.Equal(t, White, n,
					fmt.Sprintf("point (%d, %d), neighbor %d; expected White, got Black", x, y, p))
			}
		}
	}
}

func TestCountNonZeroNeighbors(t *testing.T) {
	const (
		TestCountNonZeroNeighbors1 = `
@@@
@-@
@@@
`
		TestCountNonZeroNeighbors2 = `
-@-
---
--@
`
		TestCountNonZeroNeighbors3 = `
@-@
@--
---
`
	)

	testCases := []struct {
		name     string
		bitmap   string
		expected int
		pos      image.Point
	}{
		{
			name:     "surrounded",
			bitmap:   TestCountNonZeroNeighbors1,
			expected: 8,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "2 neighbors",
			bitmap:   TestCountNonZeroNeighbors2,
			expected: 2,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "3 neighbors",
			bitmap:   TestCountNonZeroNeighbors3,
			expected: 3,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "y-axis edge",
			bitmap:   TestCountNonZeroNeighbors3,
			expected: 1,
			pos:      image.Point{X: 1, Y: 2},
		},
		{
			name:     "x-axis edge",
			bitmap:   TestCountNonZeroNeighbors3,
			expected: 0,
			pos:      image.Point{X: 2, Y: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := bitmapToImage(tc.bitmap)
			numFilled, _ := countNeighbors(img, tc.pos.X, tc.pos.Y)
			require.Equal(t, tc.expected, numFilled)
		})
	}
}

func TestCountZeroOnePatterns(t *testing.T) {
	const (
		maxCount = `
-@-
@-@
-@-
`
		maxCountOffset = `
@-@
---
@-@
`
		singleCount = `
-@-
---
---
`
		singleCountOffset1 = `
--@
---
---
`
		singleCountOffset2 = `
---
--@
---
`
		singleCountOffset3 = `
---
---
--@
`
	)

	testCases := []struct {
		name     string
		bitmap   string
		expected int
		pos      image.Point
	}{
		{
			name:     "max count",
			bitmap:   maxCount,
			expected: 4,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "max count offset",
			bitmap:   maxCountOffset,
			expected: 4,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "single count",
			bitmap:   singleCount,
			expected: 1,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "single count offset 1",
			bitmap:   singleCountOffset1,
			expected: 1,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "single count offset 2",
			bitmap:   singleCountOffset2,
			expected: 1,
			pos:      image.Point{X: 1, Y: 1},
		},
		{
			name:     "single count offset 3",
			bitmap:   singleCountOffset3,
			expected: 1,
			pos:      image.Point{X: 1, Y: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := bitmapToImage(tc.bitmap)
			actual := countZeroOnePatterns(img, 1, 1)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestThin(t *testing.T) {
	before := bitmapToImage(I)

	require.NoError(t, imgio.Imwrite(before, "before.png"))

	after, err := Thin(before)
	require.NoError(t, err)

	require.NoError(t, imgio.Imwrite(after, "after.png"))
}

const (
	DOT = `
-----
-@@@-
-@@@-
-----
`
	I = `
----------------
----@@@@@@@@----
----@@@@@@@@----
------@@@@------
------@@@@------
------@@@@------
------@@@@------
------@@@@------
----@@@@@@@@----
----@@@@@@@@----
----------------
`
	H = `
--------------------------------
---!@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@@@@@@@@@@@@@@@@@----
----@@@@@@@@@@@@@@@@@@@@@@@@----
----@@@@@@@@@@@@@@@@@@@@@@@@----
----@@@@@@@@@@@@@@@@@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
----@@@@@@@@--------@@@@@@@@----
--------------------------------
`
)

func bitmapToImage(in string) *image.Gray {
	var bm [][]string

	for _, line := range strings.Split(in, "\n") {
		if len(line) == 0 {
			// Skip empty lines (should just be the first and last lines).
			continue
		}

		var row []string

		for _, c := range strings.Split(line, "") {
			row = append(row, c)
		}

		bm = append(bm, row)
	}

	img := image.NewGray(image.Rect(0, 0, len(bm[0]), len(bm)))

	for y, row := range bm {
		for x, col := range row {
			switch col {
			case "@":
				img.Set(x, y, color.Black)
			default:
				img.Set(x, y, color.White)
			}
		}
	}

	return img
}
