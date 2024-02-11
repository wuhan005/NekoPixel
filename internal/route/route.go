package route

import (
	"github.com/flamego/flamego"
	"gorm.io/gorm"

	"github.com/wuhan005/NekoPixel/internal/context"
	"github.com/wuhan005/NekoPixel/internal/form"
)

func New(db *gorm.DB) *flamego.Flame {
	f := flamego.Classic()

	f.Use(context.Contexter(db))

	f.Group("/api", func() {
		pixelHandler := NewPixelHandler()
		f.Group("/pixels", func() {
			f.Get("/status", pixelHandler.Status)
			f.Get("/colors", pixelHandler.Colors)
			f.Combo("").
				Get(pixelHandler.GetCanvas).
				Post(form.Bind(form.SetPixels{}), pixelHandler.SetPixels)
		})
	})

	f.Get("/healthz")

	return f
}
