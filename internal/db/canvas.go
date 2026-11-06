// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"gorm.io/gorm"
)

var Canvas CanvasStore

type CanvasStore interface {
	List(ctx context.Context, fromX, fromY, toX, toY int) ([]*CanvasPixel, error)
	String(ctx context.Context, fromX, fromY, toX, toY int) (string, error)
}

func NewCanvasStore(db *gorm.DB) CanvasStore {
	return &canvas{db}
}

type CanvasPixel struct {
	UserID uint
	X      uint `gorm:"uniqueIndex:idx_canvas_pixel_x_y"`
	Y      uint `gorm:"uniqueIndex:idx_canvas_pixel_x_y"`
	Color  string
	Index  string
}

type canvas struct {
	*gorm.DB
}

func (db *canvas) List(ctx context.Context, fromX, fromY, toX, toY int) ([]*CanvasPixel, error) {
	var pixels []*CanvasPixel
	if err := db.WithContext(ctx).Where("x >= ? AND y >= ? AND x <= ? AND y <= ?", fromX, fromY, toX, toY).Find(&pixels).Error; err != nil {
		return nil, err
	}
	return pixels, nil
}

func (db *canvas) String(ctx context.Context, fromX, fromY, toX, toY int) (string, error) {
	// Leave gaps for unpainted coordinates so sparse canvases keep their row-major positions.
	const query = `
WITH positioned AS (
    SELECT (y - ?) * ? + (x - ?) AS position, index
    FROM canvas_pixels
    WHERE x >= ? AND y >= ? AND x < ? AND y < ?
), spaced AS (
    SELECT position, index, position - LAG(position, 1, -1) OVER (ORDER BY position) - 1 AS gap
    FROM positioned
)
SELECT COALESCE(STRING_AGG(REPEAT(' ', gap::integer) || index, '' ORDER BY position), '')
FROM spaced`
	var str string
	if err := db.WithContext(ctx).Raw(query, fromY, toX-fromX, fromX, fromX, fromY, toX, toY).Scan(&str).Error; err != nil {
		return "", err
	}
	return str, nil
}
