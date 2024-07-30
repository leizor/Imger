package thin

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/ernyoke/imger/imgio"
)

func TestFoo(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 100, 100))

	img.Set(30, 20, color.White)

	err := imgio.Imwrite(img, "test.png")
	if err != nil {
		t.Fatal(err)
	}

	w := color.Gray16{Y: 0xffff}
	if color.White == w {
		fmt.Println("WOOHOO")
	}
}
