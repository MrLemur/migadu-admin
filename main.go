package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
)

func main() {

	_, _, domainsString, err := loadEnvVars()
	if err != nil {
		panic(err)
	}

	domains := strings.Split(domainsString, ",")

	for _, domain := range domains {
		if !strings.Contains(domain, ".") {
			panic(fmt.Sprintf("Domain %s does not seem like a valid domain", domain))
		}
	}

	r := chi.NewRouter()

	r.Use(AddDetails)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "frontend"))
	FileServer(r, "/", filesDir)

	domainsJSON, err := json.Marshal(domains)
	if err != nil {
		panic(err)
	}

	r.Get("/api/domains", func(w http.ResponseWriter, r *http.Request) {
		w.Write(domainsJSON)
	})

	r.Mount("/api/{domain}/mailboxes", mailboxesResource{}.Routes())
	r.Mount("/api/{domain}/identities", identitiesResource{}.Routes())

	log.Fatal(http.ListenAndServe(":5000", r))
}

func AddDetails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminEmail, APIKey, _, err := loadEnvVars()
		if err != nil {
			panic(err)
		}

		ctx := context.WithValue(r.Context(), "adminEmail", adminEmail)
		ctx = context.WithValue(ctx, "APIKey", APIKey)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func loadEnvVars() (string, string, string, error) {
	godotenv.Load()

	adminEmail := os.Getenv("MIGADU_ADMIN_EMAIL")
	apiKey := os.Getenv("MIGADU_API_KEY")

	if adminEmail == "" || apiKey == "" {
		return "", "", "", fmt.Errorf("MIGADU_ADMIN_EMAIL and MIGADU_API_KEY must be set")
	}

	domainsString := os.Getenv("MIGADU_DOMAINS")

	if domainsString == "" {
		return "", "", "", fmt.Errorf("MIGADU_DOMAINS needs to have at least one domain set")
	}

	return adminEmail, apiKey, domainsString, nil
}
