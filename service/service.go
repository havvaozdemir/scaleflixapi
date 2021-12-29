package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"scaleflixapi/config"
	"scaleflixapi/data"
	types "scaleflixapi/errors"
	"scaleflixapi/logger"
	"scaleflixapi/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//Manager describes service interface
type Manager interface {
	GetMovies(resp http.ResponseWriter, req *http.Request)
	AddMovie(resp http.ResponseWriter, req *http.Request)
	GetMovieByID(resp http.ResponseWriter, req *http.Request)
	GetSeries(resp http.ResponseWriter, req *http.Request)
	AddSeries(resp http.ResponseWriter, req *http.Request)
	GetSeriesByID(resp http.ResponseWriter, req *http.Request)
	GetSuggestions(resp http.ResponseWriter, req *http.Request)
	DeleteMediaByID(resp http.ResponseWriter, req *http.Request)
	GetToken(resp http.ResponseWriter, req *http.Request)
	Authorize(next http.Handler) http.Handler
	CheckCors() http.Handler
	GetFavorites(resp http.ResponseWriter, req *http.Request)
	AddFavorite(resp http.ResponseWriter, req *http.Request)
	DeleteFavoriteByID(resp http.ResponseWriter, req *http.Request)
}

//service describes properties for api
type service struct {
	Data    data.Manager
	isAdmin bool
}

//New creates new service
func New(db *gorm.DB) Manager {
	return &service{Data: data.New(db)}
}

// swagger:route POST /movies with body
// Adds movie to database
// responses:
// 201: StatusCreated
// 400: StatusBadRequest
// 403: StatusForbidden NotAllowedAction

//AddMovie adds movie service
func (s *service) AddMovie(resp http.ResponseWriter, req *http.Request) {
	if !s.isAdmin {
		utils.WriteResponse(resp, http.StatusForbidden, types.NotAllowedAction)
		return
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resStr := fmt.Sprintf("Reading body failed: %s", err)
		logger.Error.Println(resStr)
		utils.WriteResponse(resp, http.StatusBadRequest, resStr)
		return
	}
	err = s.Data.AddMovie(body)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
	}
	utils.WriteResponse(resp, http.StatusCreated, http.StatusText(http.StatusCreated))
}

// swagger:route GET /movies movieslist
// Gets movies from database filters name and genre
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetMovies gets movies service
func (s *service) GetMovies(resp http.ResponseWriter, req *http.Request) {
	var name, genre string
	if key, ok := req.URL.Query()["name"]; ok {
		name = key[0]
	}
	if key, ok := req.URL.Query()["genre"]; ok {
		genre = key[0]
	}
	media, err := s.Data.GetMovies(name, genre)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
	}
	utils.WriteResponse(resp, http.StatusOK, media)

}

// swagger:route GET /movies/id movieslist
// Gets movie from database given id
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetMovieByID gets movie by id service
func (s *service) GetMovieByID(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key, ok := params["id"]
	if !ok {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	movie, err := s.Data.GetMovieByID(key)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, movie)
}

// swagger:route GET /series/id serie
// Gets serie from database given id
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetSeriesByID gets serie by id service
func (s *service) GetSeriesByID(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key, ok := params["id"]
	if !ok {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	series, err := s.Data.GetSeriesByID(key)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, series)
}

// swagger:route POST /series with body
// Adds serie to database
// responses:
// 201: StatusCreated
// 400: StatusBadRequest
// 403: StatusForbidden NotAllowedAction

//AddSeries adds simple series service
func (s *service) AddSeries(resp http.ResponseWriter, req *http.Request) {
	if !s.isAdmin {
		utils.WriteResponse(resp, http.StatusForbidden, types.NotAllowedAction)
		return
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resStr := fmt.Sprintf("Reading body failed: %s", err)
		logger.Error.Println(resStr)
		utils.WriteResponse(resp, http.StatusBadRequest, resStr)
		return
	}
	err = s.Data.AddSeries(body)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
	}
	utils.WriteResponse(resp, http.StatusCreated, http.StatusText(http.StatusCreated))
}

// swagger:route GET /series serieslist
// Gets series from database filters name and genre
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetSeries gets series service
func (s *service) GetSeries(resp http.ResponseWriter, req *http.Request) {
	var name, genre string
	if key, ok := req.URL.Query()["name"]; ok {
		name = key[0]
	}
	if key, ok := req.URL.Query()["genre"]; ok {
		genre = key[0]
	}

	media, err := s.Data.GetSeries(name, genre)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, media)

}

// swagger:route GET /suggestions api
// Gets fovarites from database given userId filter name and genre
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetSuggestions gets suggestions from api service
func (s *service) GetSuggestions(resp http.ResponseWriter, req *http.Request) {
	if !s.isAdmin {
		utils.WriteResponse(resp, http.StatusForbidden, types.NotAllowedAction)
		return
	}
	var name string
	if key, ok := req.URL.Query()["name"]; ok {
		name = url.QueryEscape(key[0])
	}
	if name == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	apiResp, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?t=%s&apikey=%s", name, config.APIKey))
	if utils.CheckError(err) {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	defer apiResp.Body.Close()
	body, err := ioutil.ReadAll(apiResp.Body)
	if utils.CheckError(err) {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	gelen, err := s.Data.ConvertToAPIContent(body)
	if utils.CheckError(err) {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	if gelen.Type == "series" {
		countSeason, err := strconv.Atoi(gelen.TotalSeasons)
		if err != nil {
			countSeason = 0
		}

		seasonsArr := make([]data.SeasonsAPIContent, 0, countSeason)
		for i := 0; i < countSeason; i++ {
			seasonsAPI, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?t=%s&season=%d&apikey=%s", name, i+1, config.APIKey))
			if utils.CheckError(err) {
				break
			}
			defer seasonsAPI.Body.Close()
			bodySeasons, err := ioutil.ReadAll(seasonsAPI.Body)
			if utils.CheckError(err) {
				break
			}
			seasonsAPIContent, err := s.Data.ConvertToAPISeasonsContent(bodySeasons)
			if utils.CheckError(err) {
				break
			}

			for j := 0; j < len(seasonsAPIContent.Episodes); j++ {
				episodeContentAPI, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", seasonsAPIContent.Episodes[j].ImdbID, config.APIKey))
				if utils.CheckError(err) {
					break
				}
				defer episodeContentAPI.Body.Close()
				body, err := ioutil.ReadAll(episodeContentAPI.Body)
				if utils.CheckError(err) {
					break
				}

				episodeAPIContent, err := s.Data.ConvertToAPIContent(body)
				if utils.CheckError(err) {
					break
				}

				seasonsAPIContent.Episodes[j].EpisodeContent = episodeAPIContent
			}
			seasonsArr = append(seasonsArr, seasonsAPIContent)
		}

		series := s.Data.ConvertToMedia(gelen, seasonsArr)
		utils.WriteResponse(resp, http.StatusOK, series)
		return
	}

	media := s.Data.ConvertToMedia(gelen, nil)
	utils.WriteResponse(resp, http.StatusOK, media)
}

// swagger:route DELETE /movies/{id} with body
// Deletes media from database /series/{id}
// responses:
// 200: StatusOK
// 400: StatusBadRequest
// 403: StatusForbidden NotAllowedAction

//DeleteMediaByID gets movie by id service
func (s *service) DeleteMediaByID(resp http.ResponseWriter, req *http.Request) {
	if !s.isAdmin {
		utils.WriteResponse(resp, http.StatusForbidden, types.NotAllowedAction)
		return
	}

	params := mux.Vars(req)
	key, ok := params["id"]
	if !ok {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	err := s.Data.DeleteMediaByID(key)
	if utils.CheckError(err) {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, "Succesfully deleted")
}

//GetToken gets token for given valid user information
func (s *service) GetToken(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, types.UsernamePasswordError)
		return
	}
	token, err := s.Data.GetToken(body)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, types.UsernamePasswordError)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, token)
}

//Authorize token is authorized and sets role of user
func (s *service) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Add("Vary", "Authorization")
		tokenHeader := req.Header.Get("Authorization")

		if tokenHeader == "" {
			if req.URL.String() == "/token" {
				next.ServeHTTP(resp, req)
				return
			}
			utils.WriteResponse(resp, http.StatusUnauthorized, types.NoTokenFound)
			return
		}
		headerParts := strings.Split(tokenHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteResponse(resp, http.StatusUnauthorized, types.InvalidAuthenticationTokenResponse)
			return
		}

		tokenPart := headerParts[1]

		var mySigningKey = []byte(config.SecretKey)
		token, err := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error.Println(types.TokenParseError)
				return nil, errors.New(types.TokenParseError)
			}
			return mySigningKey, nil
		})

		if err != nil {
			utils.WriteResponse(resp, http.StatusUnauthorized, types.TokenExpired)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "admin" {
				s.isAdmin = true
			} else if claims["role"] == "user" {
				s.isAdmin = false
			} else {
				utils.WriteResponse(resp, http.StatusUnauthorized, types.RoleNotImplemented)
				return
			}
		}
		next.ServeHTTP(resp, req)
	})
}

//CheckCors checks cors service
func (s *service) CheckCors() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			utils.WriteResponse(resp, http.StatusOK, "")
			return
		}
		utils.WriteResponse(resp, http.StatusMethodNotAllowed, "")
	})
}

// swagger:route GET /favorites queryparams
// Gets fovarites from database given userId filter name and genre
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//GetFavorites gets favorites for user with filter
func (s *service) GetFavorites(resp http.ResponseWriter, req *http.Request) {
	var name, genre, userID string
	if key, ok := req.URL.Query()["name"]; ok {
		name = key[0]
	}
	if key, ok := req.URL.Query()["genre"]; ok {
		genre = key[0]
	}
	if key, ok := req.URL.Query()["userId"]; ok {
		userID = key[0]
	}
	if userID == "" {
		utils.WriteResponse(resp, http.StatusBadRequest, types.UserRequired)
	}
	favorites, err := s.Data.GetFavorites(userID, name, genre)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
	}
	utils.WriteResponse(resp, http.StatusOK, favorites)
}

// swagger:route POST /favorites with jsonbody
// Adds fovarite to database
// responses:
// 201: StatusCreated
// 400: StatusBadRequest

//AddFavorite adds favorite for user given mediaid
func (s *service) AddFavorite(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resStr := fmt.Sprintf("Reading body failed: %s", err)
		logger.Error.Println(resStr)
		utils.WriteResponse(resp, http.StatusBadRequest, resStr)
		return
	}
	err = s.Data.AddFavorite(body)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
	}
	utils.WriteResponse(resp, http.StatusCreated, http.StatusText(http.StatusCreated))

}

// swagger:route DELETE /favorites/{id} with id
// Deletes fovarite from database
// responses:
// 200: StatusOK
// 400: StatusBadRequest

//DeleteFavoriteByID deletes favorites given id
func (s *service) DeleteFavoriteByID(resp http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	key, ok := params["id"]
	if !ok {
		utils.WriteResponse(resp, http.StatusBadRequest, types.KeyRequired)
		return
	}
	err := s.Data.DeleteFavoriteByID(key)
	if err != nil {
		utils.WriteResponse(resp, http.StatusBadRequest, err)
		return
	}
	utils.WriteResponse(resp, http.StatusOK, "Succesfully deleted")
}
