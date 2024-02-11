package context

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/flamego/flamego"
	"gorm.io/gorm"

	"github.com/wuhan005/NekoPixel/internal/dbutil"
)

// Context represents context of a request.
type Context struct {
	flamego.Context
}

func (c *Context) Success(data ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	}

	err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"data": d,
		},
	)
	if err != nil {
		log.Error("Failed to encode", "error", err)
	}
	return nil
}

func (c *Context) ServerError() error {
	return c.Error(http.StatusInternalServerError, "internal server error")
}

func (c *Context) Error(statusCode int, message string, v ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	c.ResponseWriter().WriteHeader(statusCode)

	err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"msg": fmt.Sprintf(message, v...),
		},
	)
	if err != nil {
		log.Error("Failed to encode", "error", err)
	}
	return nil
}

func (c *Context) Status(statusCode int) {
	c.ResponseWriter().WriteHeader(statusCode)
}

// Contexter initializes a classic context for a request.
func Contexter(gormDB *gorm.DB) flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}

		c.ResponseWriter().Header().Set("Access-Control-Allow-Origin", ctx.Request().Header.Get("Origin"))
		c.ResponseWriter().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.ResponseWriter().Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Origin, Accept, X-Requested-With")
		c.ResponseWriter().Header().Set("Access-Control-Expose-Headers", "*")
		c.ResponseWriter().Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request().Method == http.MethodOptions {
			ctx.ResponseWriter().WriteHeader(http.StatusOK)
			return
		}

		c.MapTo(gormDB, (*dbutil.Transactor)(nil))
		c.Map(c)
	}
}
