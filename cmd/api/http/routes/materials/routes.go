package materials

import "github.com/connor-davis/threereco-nextgen/internal/routing"

type Router interface {
	LoadRoutes() []routing.Route
}
