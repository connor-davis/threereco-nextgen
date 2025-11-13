package analytics

import (
	"github.com/connor-davis/threereco-nextgen/cmd/api/http/middleware"
	"github.com/connor-davis/threereco-nextgen/internal/routing"
	"github.com/connor-davis/threereco-nextgen/internal/storage"
)

type AnalyticsRouter struct {
	storage    storage.Storage
	middleware middleware.Middleware
}

func NewAnalyticsRouter(storage storage.Storage, middleware middleware.Middleware) Router {
	return &AnalyticsRouter{
		storage:    storage,
		middleware: middleware,
	}
}

func (r *AnalyticsRouter) LoadRoutes() []routing.Route {
	totalWasteCollectedRoute := r.TotalWasteCollectedRoute()

	return []routing.Route{
		totalWasteCollectedRoute,
	}
}
