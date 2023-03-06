package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *SubscriptionService) routes() http.Handler {
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
	mux.Get("/", app.Welcome)

	mux.Route("/subscription", func(mux chi.Router) {
		mux.Post("/new", app.CreateSubscription)
		mux.Get("/user/{user_id}", app.GetSubscriptionByUserID)
		mux.Get("/{id}", app.GetSubscriptionByID)

		mux.Get("/dish/{subscription_id}", app.GetDishBySubscriptionID)
		mux.Get("/delivery/{dish_id}", app.GetDishDeliveryStatus)

	})

	// mux.Get("/playlists/sort?{}", app.Playlists)
	return mux
}
