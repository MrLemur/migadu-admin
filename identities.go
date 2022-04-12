package main

import (
	"encoding/json"
	"net/http"

	"github.com/MrLemur/migadu-go"
	"github.com/go-chi/chi/v5"
)

type identitiesResource struct{}

func (rs identitiesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{mailbox}", rs.List)
	r.Post("/{mailbox}", rs.New)

	r.Route("/{mailbox}/{localPart}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs identitiesResource) List(w http.ResponseWriter, r *http.Request) {

	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	identities, err := client.ListIdentities(r.Context(), chi.URLParam(r, "mailbox"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(identities)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs identitiesResource) New(w http.ResponseWriter, r *http.Request) {
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	var expectedJSON struct {
		LocalPart   string `json:"localPart"`
		DisplayName string `json:"displayName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&expectedJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mailbox := chi.URLParam(r, "mailbox")

	if expectedJSON.LocalPart == "" || expectedJSON.DisplayName == "" {
		http.Error(w, "{'error':localPart and displayName are required'}", http.StatusBadRequest)
		return
	}

	identity, err := client.NewIdentity(r.Context(), mailbox, expectedJSON.LocalPart, expectedJSON.DisplayName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(identity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs identitiesResource) Get(w http.ResponseWriter, r *http.Request) {
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	identity, err := client.GetIdentity(r.Context(), chi.URLParam(r, "mailbox"), chi.URLParam(r, "localPart"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(identity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs identitiesResource) Update(w http.ResponseWriter, r *http.Request) {

	var identity migadu.Identity

	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	identity.LocalPart = chi.URLParam(r, "localPart")

	newIdentity, err := client.UpdateIdentity(r.Context(), chi.URLParam(r, "mailbox"), &identity)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(newIdentity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs identitiesResource) Delete(w http.ResponseWriter, r *http.Request) {
	var identity migadu.Identity
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)
	identity.LocalPart = chi.URLParam(r, "localPart")
	err := client.DeleteIdentity(r.Context(), chi.URLParam(r, "mailbox"), &identity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]string{"message": "Identity deleted"})

	w.Write(jsonResponse)
}
