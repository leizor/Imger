package threshold

import (
	"errors"
	"image"
	"image/color"

	"github.com/ernyoke/imger/histogram"
	"github.com/ernyoke/imger/utils"
)

// Method is an enum type for global threshold methods
type Method int

const (
	// ThreshBinary
	//				 _
	//				| maxVal if src(x, y) > thresh
	// dst(x, y) =	|
	//				| 0 otherwise
	//				|_
	ThreshBinary Method = iota
	// ThreshBinaryInv
	//				 _
	//				| 0 if src(x, y) > thresh
	// dst(x, y) =	|
	//				| maxVal otherwise
	//				|_
	ThreshBinaryInv
	// ThreshTrunc
	//				 _
	//				| thresh if src(x, y) > thresh
	// dst(x, y) =	|
	//				| src(x, y) otherwise
	//				|_
	ThreshTrunc
	// ThreshToZero
	//				 _
	//				| src(x, y) if src(x, y) > thresh
	// dst(x, y) =	|
	//				| 0 otherwise
	//				|_
	ThreshToZero
	// ThreshToZeroInv
	//				 _
	//				| 0 if src(x, y) > thresh
	// dst(x, y) =	|
	//				| src(x, y) otherwise
	//				|_
	ThreshToZeroInv
)

// Threshold returns a 8 bit grayscale image as result which was segmented using one of the following methods:
// ThreshBinary, ThreshBinaryInv, ThreshTrunc, ThreshToZero, ThreshToZeroInv
func Threshold(img *image.Gray, t uint8, method Method) (*image.Gray, error) {
	var transform func(color.Gray) color.Gray
	switch method {
	case ThreshBinary:
		transform = func(pixel color.Gray) color.Gray {
			if pixel.Y < t {
				return color.Gray{Y: utils.MinUint8}
			} else {
				return color.Gray{Y: utils.MaxUint8}
			}
		}
	case ThreshBinaryInv:
		transform = func(pixel color.Gray) color.Gray {
			if pixel.Y < t {
				return color.Gray{Y: utils.MaxUint8}
			} else {
				return color.Gray{Y: utils.MinUint8}
			}
		}
	case ThreshTrunc:
		transform = func(pixel color.Gray) color.Gray {
			if pixel.Y < t {
				return color.Gray{Y: pixel.Y}
			} else {
				return color.Gray{Y: t}
			}
		}
	case ThreshToZero:
		transform = func(pixel color.Gray) color.Gray {
			if pixel.Y < t {
				return color.Gray{Y: utils.MinUint8}
			} else {
				return color.Gray{Y: pixel.Y}
			}
		}
	case ThreshToZeroInv:
		transform = func(pixel color.Gray) color.Gray {
			if pixel.Y < t {
				return color.Gray{Y: pixel.Y}
			} else {
				return color.Gray{Y: utils.MinUint8}
			}
		}
	default:
		return nil, errors.New("invalid threshold method")
	}
	return threshold(img, transform), nil
}

// Threshold16 returns a grayscale image represented on 16 bits as result which was segmented using one of the following
// Methods: ThreshBinary, ThreshBinaryInv, ThreshTrunc, ThreshToZero, ThreshToZeroInv
func Threshold16(img *image.Gray16, t uint16, method Method) (*image.Gray16, error) {
	var transform func(color.Gray16) color.Gray16
	switch method {
	case ThreshBinary:
		transform = func(pixel color.Gray16) color.Gray16 {
			if pixel.Y < t {
				return color.Gray16{Y: utils.MinUint16}
			} else {
				return color.Gray16{Y: utils.MaxUint16}
			}
		}
	case ThreshBinaryInv:
		transform = func(pixel color.Gray16) color.Gray16 {
			if pixel.Y < t {
				return color.Gray16{Y: utils.MaxUint16}
			} else {
				return color.Gray16{Y: utils.MinUint16}
			}
		}
	case ThreshTrunc:
		{
			transform = func(pixel color.Gray16) color.Gray16 {
				if pixel.Y < t {
					return color.Gray16{Y: pixel.Y}
				} else {
					return color.Gray16{Y: t}
				}
			}
		}
	case ThreshToZero:
		transform = func(pixel color.Gray16) color.Gray16 {
			if pixel.Y < t {
				return color.Gray16{Y: utils.MinUint16}
			} else {
				return color.Gray16{Y: pixel.Y}
			}
		}
	case ThreshToZeroInv:
		transform = func(pixel color.Gray16) color.Gray16 {
			if pixel.Y < t {
				return color.Gray16{Y: pixel.Y}
			} else {
				return color.Gray16{Y: utils.MinUint16}
			}
		}
	default:
		return nil, errors.New("invalid threshold method")
	}
	return threshold16(img, transform), nil
}

// OtsuThreshold returns a grayscale image which was segmented using Otsu's adaptive thresholding method.
// Methods: ThreshBinary, ThreshBinaryInv, ThreshTrunc, ThreshToZero, ThreshToZeroInv
// More info about Otsu's method: https://en.wikipedia.org/wiki/Otsu%27s_method
func OtsuThreshold(img *image.Gray, method Method) (*image.Gray, error) {
	return Threshold(img, otsuThresholdValue(img), method)
}

// -------------------------------------------------------------------------------------------------------
func threshold(img *image.Gray, transform func(color.Gray) color.Gray) *image.Gray {
	size := img.Bounds().Size()
	res := image.NewGray(image.Rect(0, 0, size.X, size.Y))
	utils.ForEachGrayPixel(img, func(pixel color.Gray, x, y int) {
		res.Set(x, y, transform(pixel))
	})
	return res
}

func threshold16(img *image.Gray16, transform func(color.Gray16) color.Gray16) *image.Gray16 {
	size := img.Bounds().Size()
	res := image.NewGray16(image.Rect(0, 0, size.X, size.Y))
	utils.ForEachGray16Pixel(img, func(pixel color.Gray16, x, y int) {
		res.Set(x, y, transform(pixel))
	})
	return res
}

func otsuThresholdValue(img *image.Gray) uint8 {
	hist := histogram.HistogramGray(img)
	size := img.Bounds().Size()
	totalNumberOfPixels := size.X * size.Y

	var sumHist float64
	for i, bin := range hist {
		sumHist += float64(uint64(i) * bin)
	}

	var sumBackground float64
	var weightBackground int
	var weightForeground int

	maxVariance := 0.0
	var thresh uint8
	for i, bin := range hist {
		weightBackground += int(bin)
		if weightBackground == 0 {
			continue
		}
		weightForeground = totalNumberOfPixels - weightBackground
		if weightForeground == 0 {
			break
		}

		sumBackground += float64(uint64(i) * bin)

		meanBackground := float64(sumBackground) / float64(weightBackground)
		meanForeground := (sumHist - sumBackground) / float64(weightForeground)

		variance := float64(weightBackground) * float64(weightForeground) * (meanBackground - meanForeground) * (meanBackground - meanForeground)

		if variance > maxVariance {
			maxVariance = variance
			thresh = uint8(i)
		}
	}
	return thresh
}
