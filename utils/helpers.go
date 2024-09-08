package utils

import (
	"image"
	"image/color"
	"os"
)

func init() {
	// IteratePixels can optionally run using parallel goroutines.
	if os.Getenv("IMGER_ENABLE_PARALLEL_FOR_EACH") == "true" {
		IteratePixels = ParallelForEachPixel
	}
}

var IteratePixels = func(size image.Point, f func(x, y int)) {
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			f(x, y)
		}
	}
}

// ForEachPixel iterates through each pixel of an image and calls f, supplying the color
// and offset-compensated position of that pixel.
func ForEachPixel(img image.Image, f func(pixel color.Color, x, y int)) {
	offset := img.Bounds().Min
	IteratePixels(img.Bounds().Size(), func(x, y int) {
		f(img.At(x+offset.X, y+offset.Y), x, y)
	})
}

// ForEachGrayPixel is ForEachPixel but for image.Gray images.
func ForEachGrayPixel(img *image.Gray, f func(pixel color.Gray, x, y int)) {
	offset := img.Bounds().Min
	IteratePixels(img.Bounds().Size(), func(x, y int) {
		f(img.GrayAt(x+offset.X, y+offset.Y), x, y)
	})
}

// ForEachGray16Pixel is ForEachPixel but for image.Gray16 images.
func ForEachGray16Pixel(img *image.Gray16, f func(pixel color.Gray16, x, y int)) {
	offset := img.Bounds().Min
	IteratePixels(img.Bounds().Size(), func(x, y int) {
		f(img.Gray16At(x+offset.X, y+offset.Y), x, y)
	})
}

// ForEachRGBAPixel is ForEachPixel but for image.RGBA images.
func ForEachRGBAPixel(img *image.RGBA, f func(pixel color.RGBA, x, y int)) {
	offset := img.Bounds().Min
	IteratePixels(img.Bounds().Size(), func(x, y int) {
		f(img.RGBAAt(x+offset.X, y+offset.Y), x, y)
	})
}

// ClampInt returns min if value is lesser then min, max if value is greater them max or value if the input value is
// between min and max.
func ClampInt(value int, min int, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// ClampF64 returns min if value is lesser then min, max if value is greater them max or value if the input value is
// between min and max.
func ClampF64(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// GetMax returns the maximum value from a slice
func GetMax(v []uint64) uint64 {
	m := v[0]
	for _, value := range v {
		if m < value {
			m = value
		}
	}
	return m
}
