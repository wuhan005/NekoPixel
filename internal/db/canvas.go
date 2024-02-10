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
	var str string
	if err := db.WithContext(ctx).Raw(`SELECT STRING_AGG(t.index, '') FROM (SELECT index FROM "canvas_pixels" WHERE x >= ? AND y >= ? AND x <= ? AND y <= ? ORDER BY y, x) AS t`, fromX, fromY, toX, toY).Scan(&str).Error; err != nil {
		return "", err
	}
	return str, nil
}
