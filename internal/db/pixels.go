// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var Pixels PixelStore

type PixelStore interface {
	CreatePixel(ctx context.Context, options CreatePixelOptions) error
	GetPixel(ctx context.Context, x, y uint) (*Pixel, error)
}

func NewPixelStore(db *gorm.DB) PixelStore {
	return &pixels{db}
}

type pixels struct {
	*gorm.DB
}

type Pixel struct {
	gorm.Model
	UserID uint
	X      uint `gorm:"index:idx_pixel_x_y"`
	Y      uint `gorm:"index:idx_pixel_x_y"`
	Color  string
}

type CreatePixelOptions struct {
	UserID uint
	X, Y   uint
	Color  string
}

func (db *pixels) CreatePixel(ctx context.Context, options CreatePixelOptions) error {
	pixel := &Pixel{
		UserID: options.UserID,
		X:      options.X,
		Y:      options.Y,
		Color:  options.Color,
	}
	if err := db.WithContext(ctx).Create(pixel).Error; err != nil {
		return err
	}
	return nil
}

var ErrPixelNotFound = errors.New("pixel does not exist")

func (db *pixels) GetPixel(ctx context.Context, x, y uint) (*Pixel, error) {
	var pixel Pixel
	if err := db.WithContext(ctx).First(&pixel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPixelNotFound
		}
		return nil, err
	}
	return &pixel, nil
}
