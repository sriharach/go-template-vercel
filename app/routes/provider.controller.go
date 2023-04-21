package routes

import (
	"api-connect-mongodb-atlas/app/controllers"

	g4vercel "github.com/tbxark/g4vercel"
)

type IProviderRoute interface {
	ProviderPropsRoute(a *g4vercel.Engine)
}

type ProviderRouteTool struct {
	ProviderInterface controllers.IProviders
}

func NewProviderRoute(ac controllers.IProviders) IProviderRoute {
	return &ProviderRouteTool{
		ProviderInterface: ac,
	}
}

func (pr *ProviderRouteTool) ProviderPropsRoute(a *g4vercel.Engine) {
	group := a.Group("/api")
	group.POST("/login", pr.ProviderInterface.Login)
}
