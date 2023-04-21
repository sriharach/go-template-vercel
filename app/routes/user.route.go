package routes

import (
	"api-connect-mongodb-atlas/app/controllers"

	g4vercel "github.com/tbxark/g4vercel"
)

type IUserRoute interface {
	UserPropsRoute(a *g4vercel.Engine)
}

type UserRouteTool struct {
	UserInterface controllers.IuserController
}

func NewUserRoute(ac controllers.IuserController) IUserRoute {
	return &UserRouteTool{
		UserInterface: ac,
	}
}

func (ct *UserRouteTool) UserPropsRoute(a *g4vercel.Engine) {
	group := a.Group("/api/v1")
	group.POST("/register", ct.UserInterface.RegisterAccount)
	group.GET("/user/profile", ct.UserInterface.GetUserAccount)
	group.GET("/users", ct.UserInterface.GetUsersAccount)
}
