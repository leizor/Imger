package blend

import (
	"errors"
	"image"
	"image/color"

	"github.com/ernyoke/imger/utils"
)

// AddScalarToGray takes a grayscale image and adds an integer value to all pixels of the image. If the  result
// overflows uint8, the result will be clamped to max uint8 (255).
// Example of usage:
//
//	res := blend.AddScalarToGray(img, 56)
func AddScalarToGray(img *image.Gray, value int) *image.Gray {
	size := img.Bounds().Size()
	res := image.NewGray(image.Rect(0, 0, size.X, size.Y))

	utils.ForEachGrayPixel(img, func(pixel color.Gray, x, y int) {
		updatedPixel := int(pixel.Y) + value
		res.SetGray(x, y, color.Gray{Y: uint8(utils.ClampInt(updatedPixel, utils.MinUint8, int(utils.MaxUint8)))})
	})

	return res
}

// AddGray accepts two grayscale images and adds their pixel values. If the result for a given position overflows uint8,
// the result will be clamped to max uint8 (255).
// Example of usage:
//
//	res, err := blend.AddGray(gray1, gray2)
func AddGray(img1 *image.Gray, img2 *image.Gray) (*image.Gray, error) {
	size1 := img1.Bounds().Size()
	size2 := img2.Bounds().Size()
	if size1.X != size2.X || size1.Y != size2.Y {
		return nil, errors.New("the size of the two image does not match")
	}

	o1 := img1.Bounds().Min
	o2 := img2.Bounds().Min

	res := image.NewGray(image.Rect(0, 0, size1.X, size1.Y))
	utils.IteratePixels(size1, func(x, y int) {
		p1 := img1.GrayAt(x+o1.X, y+o1.Y)
		p2 := img2.GrayAt(x+o2.X, y+o2.Y)
		sum := utils.ClampInt(int(p1.Y)+int(p2.Y), utils.MinUint8, int(utils.MaxUint8))
		res.SetGray(x, y, color.Gray{Y: uint8(sum)})
	})

	return res, nil
}

// AddGrayWeighted accepts two grayscale images and adds their pixel values using the following equation:
// res(x, y) = img1(x, y) * w1 + img2(x, y) * w2
// If the result for a given position overflows uint8, the result will be clamped to max uint8 (255).
// If the result for a given position is negative, then it will be clamped to 0.
// Example of usage:
//
//	res, err := blend.AddGrayWeighted(gray1, 0.25, gray2, 0.75)
func AddGrayWeighted(img1 *image.Gray, w1 float64, img2 *image.Gray, w2 float64) (*image.Gray, error) {
	size1 := img1.Bounds().Size()
	size2 := img2.Bounds().Size()
	if size1.X != size2.X || size1.Y != size2.Y {
		return nil, errors.New("the size of the two image does not match")
	}

	o1 := img1.Bounds().Min
	o2 := img2.Bounds().Min

	res := image.NewGray(image.Rect(0, 0, size1.X, size1.Y))

	utils.IteratePixels(size1, func(x int, y int) {
		p1 := img1.GrayAt(x+o1.X, y+o1.Y)
		p2 := img2.GrayAt(x+o2.X, y+o2.Y)
		sum := utils.ClampF64(float64(p1.Y)*w1+float64(p2.Y)*w2, utils.MinUint8, float64(utils.MaxUint8))
		res.SetGray(x, y, color.Gray{Y: uint8(sum)})
	})

	return res, nil
}
