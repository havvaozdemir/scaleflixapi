package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"scaleflixapi/config"
	"scaleflixapi/logger"
	"scaleflixapi/service"
	"scaleflixapi/utils"
)

//DB database
var DB *gorm.DB

func handler(resp http.ResponseWriter, req *http.Request) {
	utils.WriteResponse(resp, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

//SetupDB db connection
func SetupDB(dbName string) *gorm.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, dbName)
	db, err := gorm.Open("postgres", psqlconn)
	utils.CheckError(err)
	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err = sqlDB.Ping(); err != nil {
		defer db.Close()
		logger.Fatal.Fatalf("error, not sent ping to database, %w", err)
	}

	return db
}

//CloseDB closes the database
func CloseDB() {
	DB.Close()
}

//NewServer creates scaleflix-api server with postgredb
func NewServer() {
	logger.Info.Println("db setup")
	DB = SetupDB(config.DBName)
	service := service.New(DB)

	logger.Info.Println("Server starting")

	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/movies", service.AddMovie).Methods("POST")
	r.HandleFunc("/movies", service.GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", service.GetMovieByID).Methods("GET")
	r.HandleFunc("/series", service.GetSeries).Methods("GET")
	r.HandleFunc("/series/{id}", service.GetSeriesByID).Methods("GET")
	r.HandleFunc("/series", service.AddSeries).Methods("POST")
	r.HandleFunc("/movies/{id}", service.DeleteMediaByID).Methods("DELETE")
	r.HandleFunc("/series/{id}", service.DeleteMediaByID).Methods("DELETE")
	r.HandleFunc("/suggestions", service.GetSuggestions).Methods("GET")
	r.HandleFunc("/token", service.GetToken).Methods("POST")
	r.HandleFunc("/favorites", service.AddFavorite).Methods("POST")
	r.HandleFunc("/favorites", service.GetFavorites).Methods("GET")
	r.HandleFunc("/favorites/{id}", service.DeleteFavoriteByID).Methods("DELETE")

	r.MethodNotAllowedHandler = service.CheckCors()
	logger.Info.Printf("Server started %s", config.APIPort)
	r.Use(service.Authorize)

	srv := &http.Server{
		Handler:      r,
		Addr:         config.APIPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
