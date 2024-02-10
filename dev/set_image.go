// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/wuhan005/NekoPixel/internal/db"
)

const MAX_WIDTH = 1280
const MAX_HEIGHT = 720

type Color struct {
	R, G, B uint8
}

var predefinedColors = []Color{
	{0, 0, 0},       // 0
	{255, 255, 255}, // 1
	{170, 170, 170}, // 2
	{85, 85, 85},    // 3
	{254, 211, 199}, // 4
	{255, 196, 206}, // 5
	{250, 172, 142}, // 6
	{255, 139, 131}, // 7
	{244, 67, 54},   // 8
	{233, 30, 99},   // 9
	{226, 102, 158}, // A
	{156, 39, 176},  // B
	{103, 58, 183},  // C
	{63, 81, 181},   // D
	{0, 70, 112},    // E
	{5, 113, 151},   // F
	{33, 150, 243},  // G
	{0, 188, 212},   // H
	{59, 229, 219},  // I
	{151, 253, 220}, // J
	{22, 115, 0},    // K
	{55, 169, 60},   // L
	{137, 230, 66},  // M
	{215, 255, 7},   // N
	{255, 246, 209}, // O
	{248, 203, 140}, // P
	{255, 235, 59},  // Q
	{255, 193, 7},   // R
	{255, 152, 0},   // S
	{255, 87, 34},   // T
	{184, 63, 39},   // U
	{121, 85, 72},   // V
}

func colorDistance(c1, c2 Color) float64 {
	return math.Sqrt(float64((c1.R-c2.R)*(c1.R-c2.R) + (c1.G-c2.G)*(c1.G-c2.G) + (c1.B-c2.B)*(c1.B-c2.B)))
}

func nearestColor(c Color) Color {
	var nearest Color
	minDistance := math.MaxFloat64
	for _, pc := range predefinedColors {
		d := colorDistance(c, pc)
		if d < minDistance {
			minDistance = d
			nearest = pc
		}
	}
	return nearest
}

func main() {
	imageFile := flag.String("image", "", "The image file to set.")
	flag.Parse()

	dbInstance, err := db.Init()
	if err != nil {
		panic(err)
	}

	if err := dbInstance.Exec(`TRUNCATE TABLE canvas_pixels`).Error; err != nil {
		panic(err)
	}

	file, err := os.Open(*imageFile)
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width > MAX_WIDTH {
		width = MAX_WIDTH
	}
	if height > MAX_HEIGHT {
		height = MAX_HEIGHT
	}

	canvasPixels := make([]*db.CanvasPixel, 0, MAX_WIDTH*MAX_HEIGHT)
	colorsMap := make(map[string]string)

	// 遍历每个像素点
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			color := img.At(x, y)
			r, g, b, _ := color.RGBA()
			nearestColor := nearestColor(Color{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)})
			r, g, b = uint32(nearestColor.R), uint32(nearestColor.G), uint32(nearestColor.B)
			hexColor := fmt.Sprintf("%02x%02x%02x", r, g, b)
			colorsMap[hexColor] = ""

			canvasPixels = append(canvasPixels, &db.CanvasPixel{
				UserID: 1,
				X:      uint(x),
				Y:      uint(y),
				Color:  hexColor,
			})
		}
	}

	hexes := lo.Keys(colorsMap)

	colors := make([]*db.Color, 0, len(hexes))
	if err := dbInstance.Raw(`SELECT * FROM colors WHERE color IN ?`, hexes).Find(&colors).Error; err != nil {
		panic(err)
	}
	for _, color := range colors {
		colorsMap[color.Color] = color.Index
	}

	for i := range canvasPixels {
		canvasPixels[i].Index = colorsMap[canvasPixels[i].Color]
	}

	if err := dbInstance.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(canvasPixels); i += 1000 {
			end := i + 1000
			if end > len(canvasPixels) {
				end = len(canvasPixels)
			}

			fmt.Println("Inserting", i, end)
			if err := tx.Create(canvasPixels[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
