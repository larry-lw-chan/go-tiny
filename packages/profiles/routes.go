package profiles

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/larry-lw-chan/goti/packages/sessions/flash"
	"github.com/larry-lw-chan/goti/packages/sessions/middleware"
)

// Protected Routes
func Router() http.Handler {
	r := chi.NewRouter()

	// Load Middleware
	r.Use(flash.TryGetFlash)
	r.Use(middleware.CheckIfAuthenticated)

	// Profile Routes
	r.Get("/show", ShowHandler)
	return r
}
