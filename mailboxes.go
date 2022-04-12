package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MrLemur/migadu-go"
	"github.com/go-chi/chi/v5"
)

type mailboxesResource struct{}

func (rs mailboxesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.New)

	r.Route("/{localPart}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs mailboxesResource) List(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Context())
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	mailboxes, err := client.ListMailboxes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(mailboxes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs mailboxesResource) New(w http.ResponseWriter, r *http.Request) {
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	var expectedJSON struct {
		LocalPart       string `json:"localPart"`
		DisplayName     string `json:"displayName"`
		InvitationEmail string `json:"invitationEmail"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&expectedJSON); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if expectedJSON.LocalPart == "" || expectedJSON.DisplayName == "" || expectedJSON.InvitationEmail == "" {
		http.Error(w, "{'error':'localPart, displayName and invitationEmail are required'}", http.StatusBadRequest)
		return
	}

	mailbox, err := client.NewMailbox(r.Context(), expectedJSON.LocalPart, expectedJSON.DisplayName, expectedJSON.InvitationEmail, expectedJSON.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(mailbox)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs mailboxesResource) Get(w http.ResponseWriter, r *http.Request) {
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	mailbox, err := client.GetMailbox(r.Context(), chi.URLParam(r, "localPart"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(mailbox)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs mailboxesResource) Update(w http.ResponseWriter, r *http.Request) {

	var mailbox migadu.Mailbox

	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)

	if err := json.NewDecoder(r.Body).Decode(&mailbox); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mailbox.LocalPart = chi.URLParam(r, "localPart")

	newMailbox, err := client.UpdateMailbox(r.Context(), &mailbox)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(newMailbox)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (rs mailboxesResource) Delete(w http.ResponseWriter, r *http.Request) {
	var mailbox migadu.Mailbox
	adminEmail := r.Context().Value("adminEmail").(string)
	APIKey := r.Context().Value("APIKey").(string)

	client := NewMigaduClient(chi.URLParam(r, "domain"), adminEmail, APIKey)
	mailbox.LocalPart = chi.URLParam(r, "localPart")
	err := client.DeleteMailbox(r.Context(), &mailbox)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]string{"message": "Mailbox deleted"})

	w.Write(jsonResponse)
}
