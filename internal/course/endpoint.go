package course

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonathannavas/go_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	UpdateRequest struct {
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
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

		var courseBody CreateRequest

		err := json.NewDecoder(r.Body).Decode(&courseBody)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Invalid request format"})
			return
		}

		if courseBody.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Name is required"})
			return
		}

		if courseBody.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "StartDate is required"})
			return
		}

		if courseBody.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "EndDate is required"})
			return
		}

		course, err := s.Create(courseBody.Name, courseBody.StartDate, courseBody.EndDate)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 201, Data: course})

	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		course, err := s.Get(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: course})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()
		filters := Filters{
			Name: queryParams.Get("name"),
		}

		limit, _ := strconv.Atoi(queryParams.Get("limit"))
		page, _ := strconv.Atoi(queryParams.Get("page"))

		count, err := s.Count(filters)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Error: err.Error()})
			return
		}

		meta, err := meta.New(page, limit, count)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&Response{Status: 500, Error: err.Error()})
			return
		}

		courses, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: courses, Meta: meta})
	}

}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("update users service")

		var body UpdateRequest

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Invalid request format"})
			return
		}

		if body.Name != nil && *body.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Name is required"})
			return
		}

		if body.StartDate != nil && *body.StartDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "StartDate is required"})
			return
		}

		if body.EndDate != nil && *body.EndDate == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "EndDate is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, body.Name, body.StartDate, body.EndDate); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 404, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}
