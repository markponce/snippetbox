package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/markponce/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// // Create a new instance of our application struct. For now, this just
	// // contains a structured logger (which uses the slog.DiscardHandler handler
	// // and will discard anything written to it with no action).
	// app := &application{
	// 	logger: slog.New(slog.DiscardHandler),
	// }

	// // We then use the httptest.NewTLSServer() function to create a new test
	// // server, passing in the value returned by our app.routes() method as the
	// // handler for the server. This starts up a HTTPS server which listens on a
	// // randomly-chosen port of your local machine for the duration of the test.
	// // Notice that we defer a call to ts.Close() so that the server is shutdown
	// // when the test finishes.
	// ts := httptest.NewTLSServer(app.routes())

	// defer ts.Close()

	// // The network address that the test server is listening on is contained in
	// // the ts.URL field. We can  use this along with the ts.Client().Get() method
	// // to make a GET /ping request against the test server. This returns a
	// // http.Response struct containing the response.
	// rs, err := ts.Client().Get(ts.URL + "/ping")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // We can then check the value of the response status code and body using
	// // the same pattern as before.
	// assert.Equal(t, rs.StatusCode, http.StatusOK)

	// defer rs.Body.Close()
	// body, err := io.ReadAll(rs.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// body = bytes.TrimSpace(body)

	// assert.Equal(t, string(body), "OK")

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1/",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existend ID",
			urlPath:  "/snippet/view/2/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo/",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}

}

func TestUserSignup(t *testing.T) {
	// Create the application struct containing our mocked dependencies and set
	// up the test server for running an end-to-end test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup/")
	validCSRFToken := extractCSRFToken(t, body)

	// t.Logf("CSRF token is %q", validCSRFToken)

	// _, _, body := ts.postForm(t, "/user/signup")
	// csrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup/' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup/", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	// fields:
	// name
	// content
	// expires

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create/")

		// t.Logf("code: %v", code)
		// t.Logf("coheadersde: %v", headers)
		// t.Logf("body: %v", body)

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")

	})

	t.Run("Authenticated", func(t *testing.T) {
		// Make a GET /user/login request and extract the CSRF token from the
		// response.
		_, _, body := ts.get(t, "/user/login/")
		csrfToken := extractCSRFToken(t, body)

		// Make a POST /user/login request using the extracted CSRF token and
		// credentials from our the mock user model.
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login/", form)

		// Then check that the authenticated user is shown the create snippet
		// form.
		code, _, body := ts.get(t, "/snippet/create/")

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, "<form action='/snippet/create/' method='POST'>")
	})

}
