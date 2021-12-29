// Package Classification ScaleFlix API
//
// Documentation of ScaleFlix API that represents movie and serie information with all content
//
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
// Contact: Havva Ozdemir <havvaozdemir34@gmail.com>
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"scaleflixapi/logger"
	"scaleflixapi/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal.Printf("Failed: (%v)", r)
			server.CloseDB()
		}
	}()
	server.NewServer()
}
