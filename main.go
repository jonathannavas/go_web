package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jonathannavas/go_web/internal/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	router := mux.NewRouter()

	_ = godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"))
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db = db.Debug()

	_ = db.AutoMigrate(&user.User{})

	userService := user.NewService()
	userEndpoints := user.MakeEndpoints(userService)

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users", userEndpoints.Delete).Methods("DELETE")

	srv := &http.Server{
		// Handler:      http.TimeoutHandler(router, time.Second*3, "Timeout!!"),
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
