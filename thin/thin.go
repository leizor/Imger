package thin

import (
	"fmt"
	"image"
	"image/color"
	"sync/atomic"

	"github.com/ernyoke/imger/utils"
)

// Thin returns a grayscale image which has been thinned using the algorithm detailed in
// "A Fast Parallel Algorithm for Thinning Digital Patterns" by T.Y. Zhang and C.Y. Suen.
//
// See: https://dl.acm.org/doi/pdf/10.1145/357994.358023
//
// The img provided is expected to be binary digitized picture where each pixel is either
// color.Black (1) or color.White (0).
func Thin(img *image.Gray) (*image.Gray, error) {
	res := image.NewGray(img.Bounds())
	itr := image.NewGray(img.Bounds()) // An intermediate representation; used in the main loop.

	size := img.Bounds().Size()
	offset := img.Bounds().Min

	utils.ParallelForEachPixel(size, func(x, y int) {
		res.Set(x+offset.X, y+offset.Y, img.At(x+offset.X, y+offset.Y))
		itr.Set(x+offset.X, y+offset.Y, img.At(x+offset.X, y+offset.Y))
	})

	hasDelta := atomic.Bool{}

	for {
		// Sub-iteration 1 (East/South boundary points or NW corner points)
		// Delete point p1 if it satisfies the following conditions:
		// (a) 2 <= B(P1) <= 6
		// (b) A(P1) = 1
		// (c) P2 * P4 * P6 = 0
		// (d) P4 * P6 * P8 = 0
		utils.ParallelForEachPixel(size, func(x, y int) {
			// Condition (a)
			numNeighbors := countNonZeroNeighbors(res, x+offset.X, y+offset.Y)
			aCond := 2 <= numNeighbors && numNeighbors <= 6

			// Condition (b)
			numZeroOnePatterns := countZeroOnePatterns(res, x+offset.X, y+offset.Y)
			bCond := numZeroOnePatterns == 1

			// Condition (c)
			p2 := digitize(neighbor(img, x+offset.X, y+offset.Y, 2))
			p4 := digitize(neighbor(img, x+offset.X, y+offset.Y, 4))
			p6 := digitize(neighbor(img, x+offset.X, y+offset.Y, 6))
			cCond := (p2 * p4 * p6) == 0

			// Condition (d)
			p8 := digitize(neighbor(img, x+offset.X, y+offset.Y, 8))
			dCond := (p4 * p6 * p8) == 0

			if aCond && bCond && cCond && dCond {
				itr.Set(x+offset.X, y+offset.Y, color.White)
			}
		})

		if !hasDelta.Swap(false) {
			break
		}

		copyImg(itr, res)

		// Sub-iteration2 (NW boundary points or SE corner points)
		// Delete point p1 if it satisfies the following conditions:
		// (a) 2 <= B(P1) <= 6
		// (b) A(P1) = 1
		// (c) P2 * P4 * P8 = 0
		// (d) P2 * P6 * P8 = 0
		utils.ParallelForEachPixel(size, func(x, y int) {
			// Condition (a)
			numNeighbors := countNonZeroNeighbors(res, x+offset.X, y+offset.Y)
			aCond := 2 <= numNeighbors && numNeighbors <= 6

			// Condition (b)
			numZeroOnePatterns := countZeroOnePatterns(res, x+offset.X, y+offset.Y)
			bCond := numZeroOnePatterns == 1

			// Condition (c)
			p2 := digitize(neighbor(img, x+offset.X, y+offset.Y, 2))
			p4 := digitize(neighbor(img, x+offset.X, y+offset.Y, 4))
			p8 := digitize(neighbor(img, x+offset.X, y+offset.Y, 8))
			cCond := (p2 * p4 * p8) == 0

			// Condition (d)
			p6 := digitize(neighbor(img, x+offset.X, y+offset.Y, 6))
			dCond := (p2 * p6 * p8) == 0

			if aCond && bCond && cCond && dCond {
				itr.Set(x+offset.X, y+offset.Y, color.White)
			}
		})

		if !hasDelta.Swap(false) {
			break
		}

		copyImg(itr, res)
	}

	return res, nil
}

func digitize(c color.Color) int {
	switch c {
	case color.White:
		return 0
	case color.Black:
		return 1
	default:
		// This should never happen.
		panic("unexpected color")
	}
}

// neighbor returns the color of neighboring point designation "p" to the point at (x, y)
// as follows:
//
// 9 | 2 | 3
// 8 | 1 | 4
// 7 | 6 | 5
//
// Calling neighbor with an invalid p will cause a panic.
func neighbor(img *image.Gray, x, y, p int) color.Color {
	switch p {
	case 1:
		return img.At(x, y)
	case 2:
		if y > 0 {
			return img.At(x, y-1)
		}
	case 3:
		if y > 0 && x < img.Bounds().Dx() {
			return img.At(x+1, y-1)
		}
	case 4:
		if x < img.Bounds().Dx() {
			return img.At(x+1, y)
		}
	case 5:
		if x < img.Bounds().Dx() && y < img.Bounds().Dy() {
			return img.At(x+1, y+1)
		}
	case 6:
		if y < img.Bounds().Dy() {
			return img.At(x, y+1)
		}
	case 7:
		if x > 0 && y < img.Bounds().Dy() {
			return img.At(x-1, y+1)
		}
	case 8:
		if x > 0 {
			return img.At(x-1, y)
		}
	case 9:
		if x > 0 && y > 0 {
			return img.At(x-1, y-1)
		}
	default:
		// This should never happen.
		panic(fmt.Sprintf("invalid neighbor designation p: %d", p))
	}

	// If we made it down here, that means the neighbor is off the edge of the picture;
	// return 0.
	return color.White
}

func countNonZeroNeighbors(img *image.Gray, x, y int) int {
	res := 0

	for p := 2; p < 10; p++ {
		np := neighbor(img, x, y, p)
		if np == color.Black {
			res++
		}
	}

	return res
}

// countZeroOnePatterns counts the number of 0-1 patterns in the ordered set P2, P3, P4,
// ... P8, P9.
func countZeroOnePatterns(img *image.Gray, x, y int) int {
	set := make([]color.Color, 8)
	for p := 2; p < 10; p++ {
		set[p-2] = neighbor(img, x, y, p)
	}

	res := 0
	for i, c := range set {
		if c == color.White && i+1 < len(set) && set[i+1] == color.Black {
			res++
		}
	}

	return res
}

func copyImg(src, dst *image.Gray) {
	if src.Bounds() != dst.Bounds() {
		// This should never happen.
		panic("src and dst bounds must be equal")
	}

	offset := src.Bounds().Min
	utils.ParallelForEachPixel(src.Bounds().Size(), func(x, y int) {
		dst.Set(x+offset.X, y+offset.Y, src.At(x+offset.X, y+offset.Y))
	})
}
