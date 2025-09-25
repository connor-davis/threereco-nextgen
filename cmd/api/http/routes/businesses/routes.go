package businesses

import "github.com/connor-davis/threereco-nextgen/internal/routing"

type IRouter interface {
	LoadRoutes() []routing.Route
}
