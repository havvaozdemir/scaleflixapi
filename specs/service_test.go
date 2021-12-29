package specs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scaleflixapi/config"
	"scaleflixapi/data"
	"scaleflixapi/server"
	"scaleflixapi/service"
	"testing"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
)

func initDB() *gorm.DB {
	db := server.SetupDB(config.DBNameTest)
	db.DropTableIfExists(&data.User{}, &data.Seasons{}, &data.Episodes{}, &data.Media{}, &data.UserMedia{})
	db.AutoMigrate(&data.User{}, &data.Seasons{}, &data.Episodes{}, &data.Media{}, &data.UserMedia{})
	user := CreateAdminUser()
	db.Create(&user)
	user = CreateUser()
	db.Create(&user)
	return db
}
func TestGetMovies(t *testing.T) {

	req, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := initDB()
	s := service.New(db)
	byteMovie, _ := json.Marshal(CreateTestMovie())
	data.New(db).AddMovie(byteMovie)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetMovies)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response []data.Media
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected list of movies, got: %v", rr.Body.String())
	}
	if len(response) < 1 {
		t.Errorf("expected at least 1 movie, got %v", len(response))
	}

	for _, movie := range response {
		if movie.ID == 0 {
			t.Errorf("expected movie id %d to  have a source path, was empty", movie.ID)
		}
	}
}

func TestGetSeries(t *testing.T) {
	req, err := http.NewRequest("GET", "/series", nil)
	if err != nil {
		t.Fatal(err)
	}
	db := initDB()
	s := service.New(db)
	byteSeries, _ := json.Marshal(CreateTestSeries())
	data.New(db).AddSeries(byteSeries)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetSeries)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response []data.Media
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected list of Series, got: %v", rr.Body.String())
	}
	if len(response) < 1 {
		t.Errorf("expected at least 1 Serie, got %v", len(response))
	}

	for _, serie := range response {
		if serie.ID == 0 {
			t.Errorf("expected Serie id %d to  have a source path, was empty", serie.ID)
		}
	}
}

func TestGetSeriesById(t *testing.T) {
	req, err := http.NewRequest("GET", "/series", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	db := initDB()
	s := service.New(db)
	byteSeries, _ := json.Marshal(CreateTestSeries())
	data.New(db).AddSeries(byteSeries)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetSeriesByID)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response data.Media
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected Serie, got: %v", rr.Body.String())
	}

	if response.ID != 1 {
		t.Errorf("expected Serie id %d to  have a source path, was empty", response.ID)
	}
}

func TestGetMoviesById(t *testing.T) {
	req, err := http.NewRequest("GET", "/movies", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	db := initDB()
	s := service.New(db)
	byteMovie, _ := json.Marshal(CreateTestMovie())
	data.New(db).AddMovie(byteMovie)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.GetMovieByID)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response data.Media
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("got invalid response, expected movie, got: %v", rr.Body.String())
	}

	if response.ID != 1 {
		t.Errorf("expected movie id %d to  have a source path, was empty", response.ID)
	}

}

func TestAddMovie(t *testing.T) {
	byteMovie, _ := json.Marshal(CreateTestMovie())
	reader := bytes.NewReader(byteMovie)
	db := initDB()
	s := service.New(db)
	byteUser, _ := json.Marshal(CreateAdminUser())
	token, _ := data.New(db).GetToken(byteUser)
	req, err := http.NewRequest("POST", "/movies", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := s.Authorize(http.HandlerFunc(s.AddMovie))
	req.Header.Set("Authorization", "Bearer "+token.TokenString)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAddSeries(t *testing.T) {
	byteMovie, _ := json.Marshal(CreateTestSeries())
	reader := bytes.NewReader(byteMovie)
	db := initDB()
	s := service.New(db)
	byteUser, _ := json.Marshal(CreateAdminUser())
	token, _ := data.New(db).GetToken(byteUser)
	req, err := http.NewRequest("POST", "/series", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := s.Authorize(http.HandlerFunc(s.AddSeries))
	req.Header.Set("Authorization", "Bearer "+token.TokenString)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestAddFavorite(t *testing.T) {
	byteMovie, _ := json.Marshal(CreateFavorites())
	reader := bytes.NewReader(byteMovie)
	db := initDB()
	s := service.New(db)
	req, err := http.NewRequest("POST", "/favorites", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.AddFavorite)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetSuggestions(t *testing.T) {
	testCases := map[string]struct {
		params     map[string]string
		statusCode int
	}{
		"good params": {
			map[string]string{
				"name": "Matrix",
			},
			http.StatusOK,
		},
		"without params": {
			map[string]string{},
			http.StatusBadRequest,
		},
	}
	db := initDB()
	s := service.New(db)
	byteUser, _ := json.Marshal(CreateAdminUser())
	token, _ := data.New(db).GetToken(byteUser)

	for tc, tp := range testCases {
		req, err := http.NewRequest("GET", "/suggestions", nil)
		q := req.URL.Query()
		for k, v := range tp.params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		handler := s.Authorize(http.HandlerFunc(s.GetSuggestions))
		req.Header.Set("Authorization", "Bearer "+token.TokenString)

		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != tp.statusCode {
			t.Errorf("`%v` failed, handler returned wrong status code: got %v want %v",
				tc, status, tp.statusCode)
		}
		if tc == "good params" {
			var response data.Media
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Errorf("got invalid response, expected list of movies, got: %v", rr.Body.String())
			}
			if response.Title != "Matrix" {
				t.Errorf("expected at least 1 movie named Matrix, got %v", response.Title)
			}
		}
	}
}
