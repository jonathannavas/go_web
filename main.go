package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jonathannavas/go_web/internal/course"
	"github.com/jonathannavas/go_web/internal/enrollment"
	"github.com/jonathannavas/go_web/internal/user"
	"github.com/jonathannavas/go_web/pkg/bootstrap"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()

	logs := bootstrap.InitLogger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		logs.Fatal(err)
	}

	userRepository := user.NewRepo(logs, db)
	userService := user.NewService(logs, userRepository)
	userEndpoints := user.MakeEndpoints(userService)
	usersById := "/users/{id}"

	courseRepository := course.NewRepo(logs, db)
	courseService := course.NewService(logs, courseRepository)
	courseEndpoints := course.MakeEndpoints(courseService)
	coursesById := "/courses/{id}"

	enrollmentRepository := enrollment.NewRepo(logs, db)
	enrollmentService := enrollment.NewService(logs, enrollmentRepository, userService, courseService)
	enrollmentEndpoints := enrollment.MakeEndpoints(enrollmentService)

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc(usersById, userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc(usersById, userEndpoints.Update).Methods("PATCH")
	router.HandleFunc(usersById, userEndpoints.Delete).Methods("DELETE")

	router.HandleFunc("/courses", courseEndpoints.Create).Methods("POST")
	router.HandleFunc(coursesById, courseEndpoints.Get).Methods("GET")
	router.HandleFunc("/courses", courseEndpoints.GetAll).Methods("GET")
	router.HandleFunc(coursesById, courseEndpoints.Update).Methods("PATCH")
	router.HandleFunc(coursesById, courseEndpoints.Delete).Methods("DELETE")

	router.HandleFunc("/enrollments", enrollmentEndpoints.Create).Methods("POST")

	srv := &http.Server{
		// Handler:      http.TimeoutHandler(router, time.Second*3, "Timeout!!"),
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
