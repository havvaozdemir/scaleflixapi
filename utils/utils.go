package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"scaleflixapi/logger"
)

//GetEnv gets enviroment variables
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//CheckError checks error and return true if error is not equal nil
func CheckError(err error) bool {
	if err != nil {
		logger.Error.Println(err)
		return true
	}
	return false
}

//WriteResponse writes response with given status and body
func WriteResponse(resp http.ResponseWriter, statusCode int, value interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, DELETE, POST")
	resp.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	resp.WriteHeader(statusCode)
	if err := json.NewEncoder(resp).Encode(value); err != nil {
		logger.Error.Println(err)
		return
	}

}
