package main

import (
	"dss-api/internal/data"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// TEST ADD A USER
	mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request) {
		var u = data.User{
			UserName:  "test",
			Email:     "test@test.com",
			FirstName: "You",
			LastName:  "There",
			Password:  "password",
			Level:     1,
		}

		app.infoLog.Println("Adding user..")

		id, err := app.models.User.Insert(u)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}

		app.infoLog.Println("Got back id of", id)
		newUser, err := app.models.User.GetOne(id)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			return
		}
		app.writeJSON(w, http.StatusOK, newUser)

	})

	// TEST GENERATE TOKEN
	mux.Get("/test-generate-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := app.models.Token.GenerateToken(1, 60*time.Minute)
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		token.UserName = "test"
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()

		payload := jsonResponse{
			Error:   false,
			Message: "success",
			Data:    token,
		}

		app.writeJSON(w, http.StatusOK, payload)
	})

	return mux
}
