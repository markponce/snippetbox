package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/markponce/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServer(http.Dir("./ui/static"))
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for example, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL -- any requests that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	// Unprotected application routes using the "dynamic" middleware chain.
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about/{$}", dynamic.ThenFunc(app.about))
	mux.Handle("GET /snippet/view/{id}/{$}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup/{$}", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup/{$}", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login/{$}", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login/{$}", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthetication)
	mux.Handle("GET /snippet/create/{$}", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create/{$}", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout/{$}", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /account/view/{$}", protected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update/{$}", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update/{$}", protected.ThenFunc(app.accountPasswordUpdatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
