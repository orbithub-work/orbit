package utils

import (
	"image"
	_ "image/jpeg" // Register decoders
	_ "image/png"
	"os"
)

// ComputeDHash 计算图片的差值哈希 (Difference Hash)
// 算法步骤：
// 1. 缩小尺寸为 9x8
// 2. 转为灰度
// 3. 计算差异：如果 P[x,y] > P[x+1,y] 则为 1，否则为 0
// 4. 生成 64 位哈希值
func ComputeDHash(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}
	return DHashFromImage(img)
}

// DHashFromImage 从 image.Image 对象计算 dHash
func DHashFromImage(img image.Image) (uint64, error) {
	// 1. 简单的缩放 (这里使用简单的采样，为了性能不引入额外的 resize 库)
	// 我们需要 9x8 = 72 个像素
	const width = 9
	const height = 8

	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	// 采样步长
	stepX := dx / width
	stepY := dy / height

	if stepX == 0 {
		stepX = 1
	}
	if stepY == 0 {
		stepY = 1
	}

	var pixels [height][width]uint8

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 取采样点的中心
			srcX := bounds.Min.X + x*stepX + stepX/2
			srcY := bounds.Min.Y + y*stepY + stepY/2

			// 边界检查
			if srcX >= bounds.Max.X {
				srcX = bounds.Max.X - 1
			}
			if srcY >= bounds.Max.Y {
				srcY = bounds.Max.Y - 1
			}

			r, g, b, _ := img.At(srcX, srcY).RGBA()
			// 转灰度: 0.299R + 0.587G + 0.114B
			// RGBA 返回的是 0-65535，我们需要 0-255
			gray := (299*uint32(r>>8) + 587*uint32(g>>8) + 114*uint32(b>>8)) / 1000
			pixels[y][x] = uint8(gray)
		}
	}

	// 2. 计算哈希
	var hash uint64
	for y := 0; y < height; y++ {
		for x := 0; x < width-1; x++ {
			if pixels[y][x] > pixels[y][x+1] {
				// bit index: 0 to 63
				bitIndex := uint(y*(width-1) + x)
				hash |= 1 << bitIndex
			}
		}
	}

	return hash, nil
}

// HammingDistance 计算两个哈希值的汉明距离
// 距离越小，图片越相似
// < 5: 极相似/同一张图
// > 10: 不同的图
func HammingDistance(h1, h2 uint64) int {
	xor := h1 ^ h2
	dist := 0
	for xor > 0 {
		if xor&1 == 1 {
			dist++
		}
		xor >>= 1
	}
	return dist
}
