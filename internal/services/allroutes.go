package services

import (
	"net/http"

	"github.com/leonardinius/go-service-template/internal/apiserv"
	"github.com/leonardinius/go-service-template/internal/services/version"
)

var AllRoutes = []apiserv.Route{
	newHandlerRoute(version.NewVersionServiceHandler),
}

func newHandlerRoute(pathHandler func() (string, http.Handler)) apiserv.Route {
	path, handler := pathHandler()
	return apiserv.NewRoute(path, handler)
}
