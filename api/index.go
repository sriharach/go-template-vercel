package handler

import (
	"fmt"
	"net/http"
	"os"

	"api-connect-mongodb-atlas/app/controllers"
	"api-connect-mongodb-atlas/app/routes"
	"api-connect-mongodb-atlas/pkg/configs"
	"api-connect-mongodb-atlas/pkg/middleware"

	g4vercel "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// configs.Godotenv()

	middleware.NewHandleMidleware(r)

	fmt.Println("AM_ENV:", os.Getenv("AM_ENV"))
	server := g4vercel.New()
	server.Use(g4vercel.Recovery(func(err interface{}, c *g4vercel.Context) {
		if httpError, ok := err.(g4vercel.HttpError); ok {
			c.JSON(httpError.Status, g4vercel.H{
				"message": httpError,
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, g4vercel.H{
				"messages": message,
			})
		}
	}))

	mongoAtlas := configs.InintMongodbAtlas()
	var (
		userController     = controllers.NewUserControllers(mongoAtlas, w, r)
		providerController = controllers.NewProviderControllers(mongoAtlas, w, r)

		userRoute     = routes.NewUserRoute(userController)
		providerRoute = routes.NewProviderRoute(providerController)
	)

	userRoute.UserPropsRoute(server)
	providerRoute.ProviderPropsRoute(server)

	server.GET("/", func(context *g4vercel.Context) {
		context.JSON(200, g4vercel.H{
			"message": "Hello Golang from vercel.",
		})
	})

	server.Handle(w, r)
}
