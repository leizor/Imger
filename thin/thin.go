package thin

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ernyoke/imger/imgio"
	"github.com/ernyoke/imger/utils"
)

// Thin returns a grayscale image which has been thinned using the algorithm detailed in
// "A Fast Parallel Algorithm for Thinning Digital Patterns" by T.Y. Zhang and C.Y. Suen.
//
// See: https://dl.acm.org/doi/pdf/10.1145/357994.358023
//
// The img provided is expected to be binary digitized picture where each pixel is either
// Black (1) or White (0).
func Thin(img *image.Gray) (*image.Gray, error) {
	size := img.Bounds().Size()
	offset := img.Bounds().Min

	res := image.NewGray(img.Bounds())
	utils.ForEachPixel(size, func(x, y int) {
		res.Set(x+offset.X, y+offset.Y, img.At(x+offset.X, y+offset.Y))
	})

	var toRemove []image.Point
	alternateSubIterations := true

	for i := 0; i < 20; i++ { // FIXME
		// Sub-iteration 1 (East/South boundary points or NW corner points)
		// Delete point p1 if it satisfies the following conditions:
		// (a) 2 <= B(P1) <= 6
		// (b) A(P1) = 1
		// (c) P2 * P4 * P6 = 0
		// (d) P4 * P6 * P8 = 0
		//
		// Sub-iteration 2 (North/West boundary points or SE corner points)
		// Delete point p1 if it satisfies the following conditions:
		// (a) 2 <= B(P1) <= 6
		// (b) A(P1) = 1
		// (c) P2 * P4 * P8 = 0
		// (d) P2 * P6 * P8 = 0
		utils.ForEachPixel(size, func(x, y int) {
			// Only evaluate filled positions.
			if res.At(x+offset.X, y+offset.Y) == White {
				return
			}

			numFilledNeighbors, numNeighbors := countNeighbors(res, x+offset.X, y+offset.Y)

			// Do not evaluate points that are on the edge of the image.
			if numNeighbors != 8 {
				return
			}

			// Condition (a)
			if numFilledNeighbors < 2 || numFilledNeighbors > 6 {
				return
			}

			// Condition (b)
			numZeroOnePatterns := countZeroOnePatterns(res, x+offset.X, y+offset.Y)
			if numZeroOnePatterns != 1 {
				return
			}

			if alternateSubIterations {
				// Condition (c)
				p2 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 2)
				p4 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 4)
				p6 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 6)
				if (p2 * p4 * p6) != 0 {
					return
				}

				// Condition (d)
				p8 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 8)
				if (p4 * p6 * p8) != 0 {
					return
				}
			} else {
				// Condition (c)
				p2 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 2)
				p4 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 4)
				p8 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 8)
				if (p2 * p4 * p8) != 0 {
					return
				}

				// Condition (d)
				p6 := digitizedNeighbor(img, x+offset.X, y+offset.Y, 6)
				if (p2 * p6 * p8) != 0 {
					return
				}
			}

			alternateSubIterations = !alternateSubIterations

			toRemove = append(toRemove, image.Point{
				X: x + offset.X,
				Y: y + offset.Y,
			})
		})

		if len(toRemove) == 0 {
			break
		}

		for _, pt := range toRemove {
			res.Set(pt.X, pt.Y, White)
		}

		// FIXME
		if i < 100 {
			_ = imgio.Imwrite(res, fmt.Sprintf("iter-%d.png", i))
		}

		toRemove = nil
	}

	return res, nil
}

// color.White and color.Black are 16bit but we use 8bit here.
var (
	White = color.GrayModel.Convert(color.White)
	Black = color.GrayModel.Convert(color.Black)
)

func digitizedNeighbor(img *image.Gray, x, y, p int) int {
	n, _ := neighbor(img, x, y, p)
	return digitize(n)
}

func digitize(c color.Color) int {
	switch c {
	case White:
		return 0
	case Black:
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
// Additionally, it returns true if the neighbor actually exists and false otherwise.
//
// Calling neighbor with an invalid p will cause a panic.
func neighbor(img *image.Gray, x, y, p int) (color.Color, bool) {
	switch p {
	case 1:
		return img.At(x, y), true
	case 2:
		if y > 0 {
			return img.At(x, y-1), true
		}
	case 3:
		if y > 0 && x < img.Bounds().Dx()-1 {
			return img.At(x+1, y-1), true
		}
	case 4:
		if x < img.Bounds().Dx()-1 {
			return img.At(x+1, y), true
		}
	case 5:
		if x < img.Bounds().Dx()-1 && y < img.Bounds().Dy()-1 {
			return img.At(x+1, y+1), true
		}
	case 6:
		if y < img.Bounds().Dy()-1 {
			return img.At(x, y+1), true
		}
	case 7:
		if x > 0 && y < img.Bounds().Dy()-1 {
			return img.At(x-1, y+1), true
		}
	case 8:
		if x > 0 {
			return img.At(x-1, y), true
		}
	case 9:
		if x > 0 && y > 0 {
			return img.At(x-1, y-1), true
		}
	default:
		// This should never happen.
		panic(fmt.Sprintf("invalid neighbor designation p: %d", p))
	}

	// If we made it down here, that means the neighbor is off the edge of the picture;
	// return 0.
	return White, false
}

// countNeighbors returns the number of neighbors that are filled as well as the number
// of overall neighbors (the number of neighboring pixels not off the edge of the image).
func countNeighbors(img *image.Gray, x, y int) (numFilled int, numOverall int) {
	for p := 2; p < 10; p++ {
		np, ok := neighbor(img, x, y, p)
		if np == Black {
			numFilled++
		}
		if ok {
			numOverall++
		}
	}

	return numFilled, numOverall
}

// countZeroOnePatterns counts the number of 0-1 patterns in the ordered set P2, P3, P4,
// ... P8, P9.
func countZeroOnePatterns(img *image.Gray, x, y int) int {
	set := make([]color.Color, 8)
	for p := 2; p < 10; p++ {
		set[p-2], _ = neighbor(img, x, y, p)
	}

	res := 0
	for i, c := range set {
		if i+1 < len(set) {
			if c == White && set[i+1] == Black {
				res++
			}
		} else {
			// Wrap around from P9 to P2.
			if c == White && set[0] == Black {
				res++
			}
		}
	}

	return res
}
