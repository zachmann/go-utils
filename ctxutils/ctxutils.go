package ctxutils

import (
	"github.com/gofiber/fiber/v2"

	zutils "github.com/zachmann/go-utils"
)

// FirstNonEmptyHeaderParameter checks a fiber.Ctx for multiple header
// parameters and returns the value for the first one that is set
func FirstNonEmptyHeaderParameter(c *fiber.Ctx, parameters ...string) string {
	var fncs []func() string
	for _, param := range parameters {
		fncs = append(
			fncs, func() string {
				return c.Get(param)
			},
		)
	}
	return zutils.FirstNonEmptyFnc(fncs...)
}

// FirstNonEmptyQueryParameter checks a fiber.Ctx for multiple query
// parameters and returns the value for the first one that is set
func FirstNonEmptyQueryParameter(c *fiber.Ctx, parameters ...string) string {
	var fncs []func() string
	for _, param := range parameters {
		fncs = append(
			fncs, func() string {
				return c.Query(param)
			},
		)
	}
	return zutils.FirstNonEmptyFnc(fncs...)
}
