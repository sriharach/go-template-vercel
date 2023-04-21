package middleware

import (
	"api-connect-mongodb-atlas/pkg/models"
	"api-connect-mongodb-atlas/pkg/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
	g4vercel "github.com/tbxark/g4vercel"
)

type CustomPassportInterface interface {
	DeserializeUser(c *g4vercel.Context)
	PassportJwtValidate(c *g4vercel.Context) (string, bool)
}

type HandleMidleware struct {
	r *http.Request
}

func NewHandleMidleware(r *http.Request) CustomPassportInterface {
	return &HandleMidleware{r}
}

func (h *HandleMidleware) DeserializeUser(c *g4vercel.Context) {
	str, is_jwt := h.PassportJwtValidate(c)
	if is_jwt {
		c.JSON(http.StatusUnauthorized, models.NewBaseErrorResponse(str, http.StatusUnauthorized))
		return
	}
	return
}

func (h *HandleMidleware) PassportJwtValidate(c *g4vercel.Context) (string, bool) {
	decode, _ := utils.Decode(os.Getenv("JWT_SECRET"))

	http_cookie, _ := h.r.Cookie("access_token")
	token := http_cookie.Value
	if token == "" {
		return "Unauthorized", true
	}
	_, err := jwt.ParseWithClaims(token, &utils.PayloadsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(decode), nil
	})

	if err != nil {
		return err.Error(), true
	}

	return "is_error", false

}
