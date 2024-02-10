// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/wuhan005/NekoPixel/internal/conf"
	"github.com/wuhan005/NekoPixel/internal/context"
	"github.com/wuhan005/NekoPixel/internal/db"
	"github.com/wuhan005/NekoPixel/internal/dbutil"
	"github.com/wuhan005/NekoPixel/internal/form"
)

type pixelHandler struct{}

func NewPixelHandler() *pixelHandler {
	return &pixelHandler{}
}

func (h *pixelHandler) Status(ctx context.Context) error {
	return ctx.Success(map[string]interface{}{
		"availablePixels": 1000,
	})
}

func (h *pixelHandler) Colors(ctx context.Context) error {
	return ctx.Success(conf.App.Colors)
}

func (h *pixelHandler) GetCanvas(ctx context.Context, tx dbutil.Transactor) error {
	colors, err := db.Colors.All(ctx.Request().Context())
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get colors")
		return ctx.ServerError()
	}

	colorMap := make(map[string]string, len(colors))
	for _, color := range colors {
		colorMap[color.Index] = color.Color
	}

	canvasString, err := db.Canvas.String(ctx.Request().Context(), 0, 0, conf.App.MaxWidth, conf.App.MaxHeight)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get canvas")
		return ctx.ServerError()
	}

	return ctx.Success(map[string]interface{}{
		"width":  conf.App.MaxWidth,
		"height": conf.App.MaxHeight,
		"colors": colorMap,
		"canvas": canvasString,
	})
}

func (h *pixelHandler) SetPixels(ctx context.Context, tx dbutil.Transactor, f form.SetPixels) error {
	userID, _ := strconv.Atoi(ctx.Request().Header.Get("neko-user-id"))

	for _, pixel := range f.Pixels {
		color := pixel.Color
		color = strings.ToLower(color)

		if !lo.Contains(conf.App.Colors, color) {
			return ctx.Error(http.StatusBadRequest, "颜色不在范围内")
		}

		if pixel.X >= uint(conf.App.MaxWidth) || pixel.Y >= uint(conf.App.MaxHeight) {
			return ctx.Error(http.StatusBadRequest, fmt.Sprintf("坐标 %d,%d 超出范围", pixel.X, pixel.Y))
		}
	}

	if err := tx.Transaction(func(tx *gorm.DB) error {
		pixelStores := db.NewPixelStore(tx)

		for _, pixel := range f.Pixels {
			color := strings.TrimRight(pixel.Color, "#")
			color = strings.ToLower(color)

			if err := pixelStores.CreatePixel(ctx.Request().Context(), db.CreatePixelOptions{
				UserID: uint(userID),
				X:      pixel.X,
				Y:      pixel.Y,
				Color:  color,
			}); err != nil {
				return errors.Wrap(err, "create pixel")
			}
		}
		return nil
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to set pixels")
		return ctx.ServerError()
	}

	return ctx.Success("绘制成功")
}
