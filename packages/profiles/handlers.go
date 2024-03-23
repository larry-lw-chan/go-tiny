package profiles

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/larry-lw-chan/goti/database"
	"github.com/larry-lw-chan/goti/packages/filestore"
	"github.com/larry-lw-chan/goti/packages/sessions/flash"
	"github.com/larry-lw-chan/goti/packages/users"
	"github.com/larry-lw-chan/goti/packages/utils/render"
)

/****************************************************************
* Profile Handlers
****************************************************************/

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})

	// Get user session information
	userSession := users.GetUserSession(r)
	queries := New(database.DB)

	// Get profile information or create if not exist
	profile, err := queries.GetProfileFromUserId(context.Background(), userSession.Id)
	if err != nil {
		profile, err = createProfileForUser(userSession.Id)
		handleError(w, r, err, "/profiles/show")
	}

	// Load username and profile to pass to template
	data["Username"] = userSession.Username
	data["Profile"] = profile

	// Load user into profile
	render.Template(w, data, "/profiles/show.app.tmpl")
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})

	// Get user session information
	userSession := users.GetUserSession(r)
	queries := New(database.DB)

	// Get profile information
	profile, err := queries.GetProfileFromUserId(context.Background(), userSession.Id)
	if err != nil {
		log.Println(err)
	}

	// Load username and profile to pass to template
	data["Username"] = userSession.Username
	data["Profile"] = profile

	// Load user into profile
	render.Template(w, data, "/profiles/edit.app.tmpl")
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Get user session information
	userSession := users.GetUserSession(r)

	// Assign private to (1) if checkbox is checked ("private")
	var private int64 = 0
	if r.FormValue("private") == "private" {
		private = 1
	}

	// Update profile
	name, bio, link := r.FormValue("name"), r.FormValue("bio"), r.FormValue("link")
	_, err := updateProfileForUser(userSession.Id, name, bio, link, private)

	if err != nil {
		log.Println(err)
		flash.Set(w, r, flash.ERROR, "Profile update failed")
	} else {
		flash.Set(w, r, flash.SUCCESS, "Profile successfully updated")
	}

	// Redirect to profile page
	http.Redirect(w, r, "/profiles/show", http.StatusSeeOther)
}

/****************************************************************
* Profile Photo Handlers
****************************************************************/
func EditPhotoHandler(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data").(map[string]interface{})

	// Get user session information
	userSession := users.GetUserSession(r)
	queries := New(database.DB)

	// Get profile information
	profile, err := queries.GetProfileFromUserId(context.Background(), userSession.Id)
	if err != nil {
		log.Println(err)
	}

	// Load username and profile to pass to template
	data["Username"] = userSession.Username
	data["Profile"] = profile

	render.Template(w, data, "/profiles/edit-photo.app.tmpl")
}

func EditPhotoPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get user session information
	userSession := users.GetUserSession(r)

	// Parse our multipart form, 2 << 20 specifies a maximum upload of 2 MB files
	r.ParseMultipartForm(2 << 20)

	// Determine whether file is uploaded via base64 or regular file upload
	var fileBytes []byte
	var namePattern string
	var err error

	// Choose between base64 or file upload depending on form input avatar_base64
	if r.FormValue("avatar_base64") != "" {
		Base64Image := r.FormValue("avatar_base64")
		fileBytes, namePattern, err = decodeBase64Image(Base64Image)
		handleError(w, r, err, "/profiles/edit/photo")
	} else {
		fileBytes, namePattern = handleMultipartFile(w, r)
	}

	// Populate File Struct
	fileUpload := filestore.FileUpload{
		FileBytes:   fileBytes,
		NamePattern: namePattern,
		Directory:   userSession.Username + "/avatar",
	}

	// Upload file to directory
	filepath, err := filestore.Upload(fileUpload)
	handleError(w, r, err, "/profiles/edit/photo")

	// Save file path to database
	_, err = saveFilePathToProfile(userSession.Id, filepath)

	if err != nil {
		flash.Set(w, r, flash.ERROR, "Profile update failed")
	} else {
		flash.Set(w, r, flash.SUCCESS, "Profile successfully updated")
	}

	// Redirect to profile page
	http.Redirect(w, r, "/profiles/show", http.StatusSeeOther)
}

func DeletePhotoPostHandler(w http.ResponseWriter, r *http.Request) {
	// Get avatar path
	userSession := users.GetUserSession(r)
	queries := New(database.DB)

	// Get avatar path
	filename, err := queries.GetProfileAvatarFromUserId(context.Background(), userSession.Id)
	handleError(w, r, err, "/profiles/show")

	// Delete file from filestore
	err = filestore.Delete(filename.String)
	handleError(w, r, err, "/profiles/show")

	// Fill in update profile parameters
	updateProfileParams := UpdateProfileParams{
		Avatar: sql.NullString{Valid: false},
		UserID: userSession.Id,
	}

	// Delete file path from database
	_, err = queries.UpdateProfile(context.Background(), updateProfileParams)
	handleError(w, r, err, "/profiles/show")

	// Redirect to profile page
	flash.Set(w, r, flash.SUCCESS, "Profile photo successfully deleted")
	http.Redirect(w, r, "/profiles/show", http.StatusSeeOther)
}
