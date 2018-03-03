# Imger

## Currently supported
* IO (ImreadGray, ImreadGray16, ImreadRGBA, ImreadRGBA64, Imwrite). Supported extensions: jpg, jpeg, png
* Grayscale
* Threshold (Binary, BinaryInv, Trunc, ToZero, ToZeroInv)
* Image padding (BorderConstant, BorderReplicate, BorderReflect)
* Convolution
* Blur (Average - Box, Gaussian)
* Edge detection (Sobel, Laplacian)
* Resize (Nearest Neighbour, Linear, Catmull-Rom, Lanczos)

## Merge into ```master```
```
git checkout master
git merge --no-commit --no-ff dev
```