package padding

import (
	"errors"
	"image"
	"image/color"

	"github.com/ernyoke/imger/utils"
)

// Border is an enum type for supported padding types
type Border int

const (
	// BorderConstant - xxxabcdefghxxx - where x is a black ( color.Gray{0} ) pixel
	BorderConstant Border = iota
	// BorderReplicate - aaaabcdefghhhh - replicates the nearest pixel
	BorderReplicate
	// BorderReflect - cbabcdefgfed - reflects the nearest pixel group
	BorderReflect
)

// Paddings struct holds the padding sizes for each padding
type Paddings struct {
	// PaddingLeft is the size of the left padding
	PaddingLeft int
	// PaddingRight is the size of the right padding
	PaddingRight int
	// PaddingTop is the size of the top padding
	PaddingTop int
	// PaddingBottom is the size of the bottom padding
	PaddingBottom int
}

func topPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	offset := img.Bounds().Min
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		firstPixel := img.At(x-p.PaddingLeft+offset.X, p.PaddingTop+offset.Y)
		for y := 0; y < p.PaddingTop; y++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func bottomPaddingReplicate(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	offset := img.Bounds().Min
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		lastPixel := img.At(x-p.PaddingLeft+offset.X, originalSize.Y-1+offset.Y)
		for y := p.PaddingTop + originalSize.Y; y < originalSize.Y+p.PaddingTop+p.PaddingBottom; y++ {
			setPixel(x, y, lastPixel)
		}
	}
}

func leftPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		firstPixel := padded.At(p.PaddingLeft, y)
		for x := 0; x < p.PaddingLeft; x++ {
			setPixel(x, y, firstPixel)
		}
	}
}

func rightPaddingReplicate(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		lastPixel := padded.At(originalSize.X+p.PaddingLeft-1, y)
		for x := originalSize.X + p.PaddingLeft; x < originalSize.X+p.PaddingLeft+p.PaddingRight; x++ {
			setPixel(x, y, lastPixel)
		}
	}
}

func topPaddingReflect(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	offset := img.Bounds().Min
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		for y := 0; y < p.PaddingTop; y++ {
			pixel := img.At(x-p.PaddingLeft+offset.X, p.PaddingTop-y+offset.Y)
			setPixel(x, y, pixel)
		}
	}
}

func bottomPaddingReflect(img image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	offset := img.Bounds().Min
	for x := p.PaddingLeft; x < originalSize.X+p.PaddingLeft; x++ {
		for y := p.PaddingTop + originalSize.Y; y < originalSize.Y+p.PaddingTop+p.PaddingBottom; y++ {
			pixel := img.At(x-p.PaddingLeft+offset.X, originalSize.Y-(y-p.PaddingTop-originalSize.Y)-2+offset.Y)
			setPixel(x, y, pixel)
		}
	}
}

func leftPaddingReflect(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		for x := 0; x < p.PaddingLeft; x++ {
			pixel := padded.At(2*p.PaddingLeft-x, y)
			setPixel(x, y, pixel)
		}
	}
}

func rightPaddingReflect(img image.Image, padded image.Image, p Paddings, setPixel func(int, int, color.Color)) {
	originalSize := img.Bounds().Size()
	for y := 0; y < originalSize.Y+p.PaddingBottom+p.PaddingTop; y++ {
		for x := originalSize.X + p.PaddingLeft; x < originalSize.X+p.PaddingLeft+p.PaddingRight; x++ {
			pixel := padded.At(originalSize.X+p.PaddingLeft-(x-originalSize.X-p.PaddingLeft)-2, y)
			setPixel(x, y, pixel)
		}
	}
}

// PaddingGray appends padding to a given grayscale image. The size of the padding is calculated from the kernel size
// and the anchor point. Supported border types are: BorderConstant, BorderReplicate, BorderReflect.
// Example of usage:
//
//	res, err := padding.PaddingGray(img, {5, 5}, {1, 1}, BorderReflect)
//
// Note: this will add a 1px padding for the top and left borders of the image and a 3px padding fot the bottom and
// right borders of the image.
func PaddingGray(img *image.Gray, kernelSize image.Point, anchor image.Point, border Border) (*image.Gray, error) {
	p, err := calculatePaddings(kernelSize, anchor)
	if err != nil {
		return nil, err
	}
	padded := blankPaddingGray(img, p)

	switch border {
	case BorderConstant:
		// do nothing
	case BorderReplicate:
		topPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	case BorderReflect:
		topPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	default:
		return nil, errors.New("unknown border type")
	}
	return padded, nil
}

// PaddingRGBA appends padding to a given RGBA image. The size of the padding is calculated from the kernel size
// and the anchor point. Supported border types are: BorderConstant, BorderReplicate, BorderReflect.
// Example of usage:
//
//	res, err := padding.PaddingRGBA(img, {5, 5}, {1, 1}, BorderReflect)
//
// Note: this will add a 1px padding for the top and left borders of the image and a 3px padding fot the bottom and
// right borders of the image.
func PaddingRGBA(img *image.RGBA, kernelSize image.Point, anchor image.Point, border Border) (*image.RGBA, error) {
	p, err := calculatePaddings(kernelSize, anchor)
	if err != nil {
		return nil, err
	}
	padded := blankPaddingRGBA(img, p)

	switch border {
	case BorderConstant:
		// do nothing
	case BorderReplicate:
		topPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReplicate(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReplicate(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	case BorderReflect:
		topPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		bottomPaddingReflect(img, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		leftPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
		rightPaddingReflect(img, padded, p, func(x int, y int, pixel color.Color) {
			padded.Set(x, y, pixel)
		})
	default:
		return nil, errors.New("unknown border type")
	}
	return padded, nil
}

func blankPaddingGray(img *image.Gray, padding Paddings) *image.Gray {
	origSize := img.Bounds().Size()
	res := image.NewGray(image.Rect(
		0, 0,
		origSize.X+padding.PaddingLeft+padding.PaddingRight, origSize.Y+padding.PaddingTop+padding.PaddingBottom,
	))

	utils.ForEachPixel(img, func(pixel color.Color, x, y int) {
		res.Set(x+padding.PaddingLeft, y+padding.PaddingTop, pixel)
	})

	return res
}

func blankPaddingRGBA(img *image.RGBA, padding Paddings) *image.RGBA {
	origSize := img.Bounds().Size()
	res := image.NewRGBA(image.Rect(
		0, 0,
		origSize.X+padding.PaddingLeft+padding.PaddingRight, origSize.Y+padding.PaddingTop+padding.PaddingBottom,
	))

	utils.ForEachPixel(img, func(pixel color.Color, x, y int) {
		res.Set(x+padding.PaddingLeft, y+padding.PaddingTop, pixel)
	})

	return res
}

// -------------------------------------------------------------------------------------------------------
func calculatePaddings(kernelSize image.Point, anchor image.Point) (Paddings, error) {
	var p Paddings
	if kernelSize.X < 0 || kernelSize.Y < 0 {
		return p, errors.New("negative size")
	}
	if anchor.X < 0 || anchor.Y < 0 {
		return p, errors.New("negative anchor value")
	}
	if anchor.X > kernelSize.X || anchor.Y > kernelSize.Y {
		return p, errors.New("anc" + "hor value outside of the kernel")
	}

	p = Paddings{PaddingLeft: anchor.X, PaddingRight: kernelSize.X - anchor.X - 1, PaddingTop: anchor.Y, PaddingBottom: kernelSize.Y - anchor.Y - 1}

	return p, nil
}
