// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"context"

	"gorm.io/gorm"
)

var Colors ColorsStore

type ColorsStore interface {
	All(ctx context.Context) ([]*Color, error)
	IsEmpty(ctx context.Context) (bool, error)
	Create(ctx context.Context, options CreateColorOptions) error
	Truncate(ctx context.Context) error
}

func NewColorsStore(db *gorm.DB) ColorsStore {
	return &colors{db}
}

type Color struct {
	Color string
	Index string
}

type colors struct {
	*gorm.DB
}

func (db *colors) All(ctx context.Context) ([]*Color, error) {
	var colors []*Color
	if err := db.WithContext(ctx).Find(&colors).Error; err != nil {
		return nil, err
	}
	return colors, nil
}

func (db *colors) IsEmpty(ctx context.Context) (bool, error) {
	var count int64
	if err := db.WithContext(ctx).Model(&Color{}).Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

type CreateColorOptions struct {
	Colors map[string]string
}

func (db *colors) Create(ctx context.Context, options CreateColorOptions) error {
	if len(options.Colors) == 0 {
		return nil
	}

	colors := make([]*Color, 0, len(options.Colors))
	for color, index := range options.Colors {
		colors = append(colors, &Color{
			Color: color,
			Index: index,
		})
	}

	return db.WithContext(ctx).Create(colors).Error
}

func (db *colors) Truncate(ctx context.Context) error {
	return db.WithContext(ctx).Exec("TRUNCATE TABLE colors").Error
}
