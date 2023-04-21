package internal

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HttpCookieInterface interface {
	CookieFromFiber(c *fiber.Ctx, e *DefaultParamCookie) error
}

// type CtxFiber struct {
// 	c *fiber.Ctx
// }

type DefaultParamCookie struct {
	Name    string
	Value   string
	Expires int    `default:"30"`
	Path    string `default:"/"`
}

// func NewLauageCookie(c *fiber.Ctx) CtxFiber {
// 	return CtxFiber{c}
// }

func HttpCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(30 * time.Minute)
	cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
	http.SetCookie(w, &cookie)
}

func CookieFromFiber(c *fiber.Ctx, e *DefaultParamCookie) error {
	time := time.Now().Add(time.Duration(e.Expires) * time.Minute)
	c.Cookie(&fiber.Cookie{
		Name:    e.Name,
		Value:   e.Value,
		Expires: time,
		Path:    e.Path,
	})

	return c.Next()
}
