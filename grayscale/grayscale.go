package grayscale

import (
	"image"
	"image/color"

	"github.com/ernyoke/imger/utils"
)

// Grayscale takes an image on any type and returns the equivalent grayscale image represented on 8 bits.
func Grayscale(img image.Image) *image.Gray {
	size := img.Bounds().Size()
	res := image.NewGray(image.Rect(0, 0, size.X, size.Y))
	utils.ForEachPixel(img, func(pixel color.Color, x, y int) {
		res.Set(x, y, color.GrayModel.Convert(pixel))
	})
	return res
}

// Grayscale16 takes an image on any type and returns the equivalent grayscale image represented on 16 bits.
func Grayscale16(img image.Image) *image.Gray16 {
	size := img.Bounds().Size()
	res := image.NewGray16(image.Rect(0, 0, size.X, size.Y))
	utils.ForEachPixel(img, func(pixel color.Color, x, y int) {
		res.Set(x, y, color.Gray16Model.Convert(pixel))
	})
	return res
}
