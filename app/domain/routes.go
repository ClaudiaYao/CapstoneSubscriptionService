package domain

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (service *SubscriptionService) Routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	// mux.Route("/playlists", func(r chi.Router) {
	// 	r.Get("/", app.Playlists)
	// 	r.Post("/new", app.CreatePlaylist)

	// 	r.Get("/{code}", h.internalPlan.Get)
	// 	r.Put("/{code}", h.internalPlan.Update)
	// })
	mux.Get("/", service.Welcome)

	mux.Route("/subscription", func(mux chi.Router) {
		mux.Post("/new", service.CreateSubscription)
		mux.Get("/user/{user_id}", service.GetSubscriptionByUserID)
		mux.Get("/{id}", service.GetSubscriptionByID)

		mux.Get("/dish/{subscription_id}", service.GetDishBySubscriptionID)
		mux.Get("/delivery/{dish_id}", service.GetDishDeliveryStatus)

	})

	// mux.Get("/playlists/sort?{}", app.Playlists)
	return mux
}
