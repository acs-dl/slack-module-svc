package api

import (
	"fmt"

	"github.com/acs-dl/slack-module-svc/internal/data"
	"github.com/acs-dl/slack-module-svc/internal/data/postgres"
	"github.com/acs-dl/slack-module-svc/internal/service/api/handlers"
	"github.com/acs-dl/slack-module-svc/internal/service/api/middleware"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (r *router) apiRouter() chi.Router {
	router := chi.NewRouter()
	logger := r.cfg.Log().WithField("service", fmt.Sprintf("%s-api", data.ModuleName))
	secret := r.cfg.JwtParams().Secret

	router.Use(
		ape.RecoverMiddleware(logger),
		ape.LoganMiddleware(logger),
		ape.CtxMiddleware(
			//base
			handlers.CtxLog(logger),

			// storage
			handlers.CtxPermissionsQ(postgres.NewPermissionsQ(r.cfg.DB())),
			handlers.CtxUsersQ(postgres.NewUsersQ(r.cfg.DB())),
			handlers.CtxConversationsQ(postgres.NewConversationsQ(r.cfg.DB())),
			handlers.CtxLinksQ(postgres.NewLinksQ(r.cfg.DB())),

			// other configs
			handlers.CtxParentContext(r.ctx),
		),
	)

	router.Route("/integrations/slack", func(r chi.Router) {
		r.Get("/role", handlers.GetRole)               // comes from orchestrator
		r.Get("/roles", handlers.GetRolesMap)          // comes from orchestrator
		r.Get("/user_roles", handlers.GetUserRolesMap) // comes from orchestrator

		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", handlers.GetUserById) // comes from orchestrator
			r.Get("/", handlers.GetUsers)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.IsAuthenticated(secret))

			r.Get("/get_available_roles", handlers.GetRoles)
			r.Get("/submodule", handlers.CheckSubmodule)
			r.Get("/permissions", handlers.GetPermissions)

			r.Route("/estimate_refresh", func(r chi.Router) {
				r.Post("/submodule", handlers.GetEstimatedRefreshSubmodule)
				r.Post("/module", handlers.GetEstimatedRefreshModule)
			})
		})

		r.With(middleware.IsAdmin(secret)).Route("/links", func(r chi.Router) {
			r.Post("/", handlers.AddLink)
			r.Delete("/", handlers.RemoveLink)
		})
	})

	return router
}
