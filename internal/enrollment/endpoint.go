package enrollment

import (
	"encoding/json"
	"net/http"

	"github.com/jonathannavas/go_web/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	CreateRequest struct {
		UserID   string `json:"user_id"`
		CourseID string `json:"course_id"`
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
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var enrollmentBody CreateRequest

		err := json.NewDecoder(r.Body).Decode(&enrollmentBody)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Invalid request format"})
			return
		}

		if enrollmentBody.UserID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "User id is required"})
			return
		}

		if enrollmentBody.CourseID == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: "Course id is required"})
			return
		}

		enrollment, err := s.Create(enrollmentBody.UserID, enrollmentBody.CourseID)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Error: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 201, Data: enrollment})

	}
}
