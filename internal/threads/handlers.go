package threads

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/larry-lw-chan/goti/database"
	"github.com/larry-lw-chan/goti/internal/auth"
	"github.com/larry-lw-chan/goti/internal/sessions/flash"
	"github.com/larry-lw-chan/goti/internal/utils/render"
)

/********************************************************
* Threads
*********************************************************/
func NewThreadHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})
	render.Template(w, data, "/threads/new.app.tmpl")
}

func NewPostThreadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Handle Form Validation
	errs := validateNewThread(r)
	if errs != nil {
		var message string
		for _, err := range errs {
			message += err.Error() + "\n"
		}
		flash.Set(w, r, flash.ERROR, message)
		http.Redirect(w, r, "/threads/new", http.StatusSeeOther)
		return
	}

	// Get user session information
	userSession := auth.GetUserSession(r)

	// Create new thread
	queries := New(database.DB)
	threadParam := CreateThreadParams{
		Content:   r.FormValue("content"),
		ThreadID:  sql.NullInt64{Valid: false},
		UserID:    userSession.ID,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}
	_, err := queries.CreateThread(r.Context(), threadParam)
	if err != nil {
		flash.Set(w, r, flash.ERROR, "Failed to create new thread.  Please contact support.")
		http.Redirect(w, r, "/threads/new", http.StatusSeeOther)
		return
	}

	// Set flash message and render template
	flash.Set(w, r, flash.SUCCESS, "New thread created.")
	http.Redirect(w, r, "/threads/new", http.StatusSeeOther)
}

func ShowThreadHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})

	// Get thread_id from URL
	thread_id := chi.URLParam(r, "thread_id")

	log.Println(thread_id)

	render.Template(w, data, "/threads/show.app.tmpl")
}

/********************************************************
* Partials
*********************************************************/
func UserThreadsHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})

	// Get user session information
	userSession := auth.GetUserSession(r)

	// Temp Solution - Get all Threads
	queries := New(database.DB)
	threads, err := queries.GetUserThreads(r.Context(), userSession.ID)

	if err != nil {
		// Handle Error
		// flash.Set(w, r, flash.ERROR, "Failed to get threads.  Please contact support.")
		data["Threads"] = nil
	} else {
		data["Threads"] = threads
	}

	render.Partial(w, data, "/threads/partials/user_threads.app.tmpl", "/partials/thread.app.tmpl")
}

func UserRepliesHandler(w http.ResponseWriter, r *http.Request) {
	render.Partial(w, nil, "/threads/partials/user_replies.app.tmpl")
}

func UserRepostHandler(w http.ResponseWriter, r *http.Request) {
	render.Partial(w, nil, "/threads/partials/user_repost.app.tmpl")
}
