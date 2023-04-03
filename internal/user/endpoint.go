package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		var userBody CreateRequest

		err := json.NewDecoder(r.Body).Decode(&userBody)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"Invalid request"})
			return
		}

		if userBody.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"FirstName is required"})
			return
		}

		if userBody.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"LastName is required"})
			return
		}

		user, err := s.Create(userBody.FirstName, userBody.LastName, userBody.Email, userBody.Phone)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"LastName is required"})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		user, err := s.Get(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.GetAll()
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("update users")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("delete users")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}
