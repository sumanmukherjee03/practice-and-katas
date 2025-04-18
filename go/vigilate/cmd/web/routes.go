package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tsawler/vigilate/internal/handlers"
)

func routes() http.Handler {

	mux := chi.NewRouter()

	// default middleware
	mux.Use(SessionLoad)
	mux.Use(RecoverPanic)
	mux.Use(NoSurf)
	mux.Use(CheckRemember)

	// login
	mux.Get("/", handlers.Repo.LoginScreen)
	mux.Post("/", handlers.Repo.Login)

	mux.Get("/user/logout", handlers.Repo.Logout)

	// pusher routes
	mux.Route("/pusher", func(mux chi.Router) {
		// all nested pusher routes are protected
		mux.Use(Auth)
		mux.Post("/auth", handlers.Repo.PusherAuth)
	})

	// admin routes
	mux.Route("/admin", func(mux chi.Router) {
		// all admin routes are protected
		mux.Use(Auth)

		// sample route to send message to a private channel
		mux.Get("/private-message", handlers.Repo.SendPrivateMessage)

		// overview
		mux.Get("/overview", handlers.Repo.AdminDashboard)

		// events
		mux.Get("/events", handlers.Repo.Events)

		// settings
		mux.Get("/settings", handlers.Repo.Settings)
		mux.Post("/settings", handlers.Repo.PostSettings)

		// service status pages (all hosts)
		mux.Get("/all-healthy", handlers.Repo.AllHealthyServices)
		mux.Get("/all-warning", handlers.Repo.AllWarningServices)
		mux.Get("/all-problems", handlers.Repo.AllProblemServices)
		mux.Get("/all-pending", handlers.Repo.AllPendingServices)

		// users
		mux.Get("/users", handlers.Repo.AllUsers)
		mux.Get("/user/{id}", handlers.Repo.OneUser)
		mux.Post("/user/{id}", handlers.Repo.PostOneUser)
		mux.Get("/user/delete/{id}", handlers.Repo.DeleteUser)

		// schedule
		mux.Get("/schedule", handlers.Repo.ListEntries)

		// preference
		mux.Post("/preference/ajax/set-system-pref", handlers.Repo.SetSystemPref)
		mux.Post("/preference/ajax/toggle-monitoring", handlers.Repo.ToggleMonitoring)

		// hosts
		mux.Get("/host/all", handlers.Repo.AllHosts)
		mux.Get("/host/{id}", handlers.Repo.Host)
		mux.Post("/host/{id}", handlers.Repo.PostHost)
		mux.Post("/host/{id}/service/{service_id}/toggle", handlers.Repo.ToggleServiceForHost)
		mux.Post("/host/{id}/service/{service_id}/perform-check", handlers.Repo.PerformCheckOnServiceForHost)
	})

	// static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
