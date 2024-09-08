package edgedetection

import (
	"image"
	"testing"

	"github.com/ernyoke/imger/imgio"
	"github.com/ernyoke/imger/padding"
)

// -----------------------------Acceptance tests------------------------------------

func setupTestCaseGraySobel(t *testing.T) *image.Gray {
	path := "../res/engine.png"
	img, err := imgio.ImreadGray(path)
	if err != nil {
		t.Errorf("Could not read image from path: %s", path)
	}
	return img
}

func setupTestCaseRGBASobel(t *testing.T) *image.RGBA {
	path := "../res/engine.png"
	img, err := imgio.ImreadRGBA(path)
	if err != nil {
		t.Errorf("Could not read image from path: %s", path)
	}
	return img
}

func tearDownTestCaseSobel(t *testing.T, img image.Image, path string) {
	err := imgio.Imwrite(img, path)
	if err != nil {
		t.Errorf("Could not write image to path: %s", path)
	}
}

func Test_Acceptance_HorizontalSobelGray(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	sobel, _ := HorizontalSobelGray(gray, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/horizontalSobelGray.png")
}

func Test_Acceptance_HorizontalSobelGray_Cropped(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	cropped := gray.SubImage(image.Rect(100, 100, gray.Bounds().Size().X-100, gray.Bounds().Size().Y-100)).(*image.Gray)
	sobel, _ := HorizontalSobelGray(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/horizontalSobelGrayCropped.png")
}

func Test_Acceptance_VerticalSobelGray(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	sobel, _ := VerticalSobelGray(gray, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/verticalSobelGray.png")
}

func Test_Acceptance_VerticalSobelGray_Cropped(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	cropped := gray.SubImage(image.Rect(100, 100, gray.Bounds().Size().X-100, gray.Bounds().Size().Y-100)).(*image.Gray)
	sobel, _ := VerticalSobelGray(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/verticalSobelGrayCropped.png")
}

func Test_Acceptance_SobelGray(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	sobel, _ := SobelGray(gray, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/sobelGray.png")
}

func Test_Acceptance_SobelGray_Cropped(t *testing.T) {
	gray := setupTestCaseGraySobel(t)
	cropped := gray.SubImage(image.Rect(100, 100, gray.Bounds().Size().X-100, gray.Bounds().Size().Y-100)).(*image.Gray)
	sobel, _ := SobelGray(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/sobelGrayCropped.png")
}

func Test_Acceptance_HorizontalSobelRGBA(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	sobel, _ := HorizontalSobelRGBA(rgba, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/horizontalSobelRGBA.png")
}

func Test_Acceptance_HorizontalSobelRGBA_Cropped(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	cropped := rgba.SubImage(image.Rect(100, 100, rgba.Bounds().Size().X-100, rgba.Bounds().Size().Y-100)).(*image.RGBA)
	sobel, _ := HorizontalSobelRGBA(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/horizontalSobelRGBACropped.png")
}

func Test_Acceptance_VerticalSobelRGBA(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	sobel, _ := VerticalSobelRGBA(rgba, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/verticalSobelRGBA.png")
}

func Test_Acceptance_VerticalSobelRGBA_Cropped(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	cropped := rgba.SubImage(image.Rect(100, 100, rgba.Bounds().Size().X-100, rgba.Bounds().Size().Y-100)).(*image.RGBA)
	sobel, _ := VerticalSobelRGBA(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/verticalSobelRGBACropped.png")
}

func Test_Acceptance_SobelRGBA(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	sobel, _ := SobelRGBA(rgba, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/sobelRGBA.png")
}

func Test_Acceptance_SobelRGBA_Cropped(t *testing.T) {
	rgba := setupTestCaseRGBASobel(t)
	cropped := rgba.SubImage(image.Rect(100, 100, rgba.Bounds().Size().X-100, rgba.Bounds().Size().Y-100)).(*image.RGBA)
	sobel, _ := SobelRGBA(cropped, padding.BorderReflect)
	tearDownTestCaseSobel(t, sobel, "../res/edge/sobelRGBACropped.png")
}

// ---------------------------------------------------------------------------------
