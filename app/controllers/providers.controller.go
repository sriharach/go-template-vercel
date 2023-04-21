package controllers

import (
	"api-connect-mongodb-atlas/pkg/models"
	"api-connect-mongodb-atlas/pkg/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	g4vercel "github.com/tbxark/g4vercel"
)

type IProviders interface {
	Login(c *g4vercel.Context)
}

type PropsProviderController struct {
	MongoDB            *mongo.Database
	MainCollectionDB   *mongo.Collection
	HttpResponseWriter http.ResponseWriter
	HttpRequest        *http.Request
}

func NewProviderControllers(DB *mongo.Database, w http.ResponseWriter, r *http.Request) IProviders {
	return &PropsProviderController{
		MongoDB:            DB,
		MainCollectionDB:   DB.Collection("users"),
		HttpResponseWriter: w,
		HttpRequest:        r,
	}
}

func (pv *PropsProviderController) Login(c *g4vercel.Context) {
	var result *models.ModuleProfile
	payload := new(models.SignInInput)
	collection := pv.MainCollectionDB

	if err := json.NewDecoder(c.Req.Body).Decode(payload); err != nil {
		c.JSON(http.StatusBadRequest, models.NewBaseErrorResponse(err.Error(), http.StatusBadRequest))
	}

	bson := bson.M{
		"e_mail": payload.E_mail,
	}
	err := collection.FindOne(context.Background(), bson).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewBaseErrorResponse(err.Error(), fiber.StatusNotFound))
		return
	}

	if is_passwor_hash := utils.CheckPasswordHash(payload.Password, result.Password); !is_passwor_hash {
		c.JSON(http.StatusNotAcceptable, models.NewBaseErrorResponse("Password don't matching", fiber.StatusNotAcceptable))
		return
	}

	access_token, _ := utils.GenerateTokenJWT(result, true)
	refresh_token, _ := utils.GenerateTokenJWT(result, false)

	decode, _ := utils.Decode(os.Getenv("JWT_SECRET"))
	token, _err := jwt.ParseWithClaims(access_token, &utils.PayloadsClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(decode), nil
	})

	if _err != nil {
		c.JSON(http.StatusBadRequest, models.NewBaseErrorResponse(_err.Error(), fiber.StatusBadRequest))
	}

	claims := token.Claims.(*utils.PayloadsClaims)

	expiration := time.Now().UTC().Add(30 * time.Minute).UTC()
	cookie_access_token := http.Cookie{Name: "access_token", Value: access_token, Expires: expiration}
	cookie_user_id := http.Cookie{Name: "user_id", Value: result.ID.Hex(), Expires: expiration}
	http.SetCookie(pv.HttpResponseWriter, &cookie_access_token)
	http.SetCookie(pv.HttpResponseWriter, &cookie_user_id)

	c.JSON(http.StatusOK, models.NewBaseResponse(utils.GenerateJWTOption{
		Access_token:  access_token,
		Refresh_token: refresh_token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.IssuedAt,
		},
	}, http.StatusOK))
}
