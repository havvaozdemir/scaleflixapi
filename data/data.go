package data

import (
	"encoding/json"
	"errors"
	"scaleflixapi/config"
	types "scaleflixapi/errors"
	"scaleflixapi/logger"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"

	//all time imports
	_ "github.com/lib/pq"
)

//Manager interface for data
type Manager interface {
	AddMovie([]byte) error
	GetMovies(name, genre string) ([]Media, error)
	GetMovieByID(id string) (Media, error)
	AddSeries([]byte) error
	DeleteMediaByID(key string) error
	GetSeries(name, genre string) ([]Media, error)
	GetSeriesByID(id string) (Media, error)
	ConvertToMedia(fromAPIContent MediaAPIContent, fromAPISeasons []SeasonsAPIContent) *Media
	ConvertToAPIContent(body []byte) (MediaAPIContent, error)
	ConvertToAPISeasonsContent(body []byte) (SeasonsAPIContent, error)
	GetToken(body []byte) (Token, error)
	AddFavorite([]byte) error
	DeleteFavoriteByID(key string) error
	GetFavorites(userID, name, genre string) ([]UserMedia, error)
}

//MediaType definition
type MediaType int

const (
	//Movie enum
	Movie MediaType = iota + 1
	//Series enum
	Series
	//Episode enum
	Episode
)

//Data definition
type Data struct {
	DB *gorm.DB
}

//MediaAPIContent definition
type MediaAPIContent struct {
	Type         string `json:"Type"`
	Title        string `json:"Title"`
	Plot         string `json:"Plot"`
	ImdbRating   string `json:"imdbRating"`
	Director     string `json:"Director"`
	Writer       string `json:"Writer"`
	Aktors       string `json:"Aktors"`
	Released     string `json:"Released"`
	Runtime      string `json:"Runtime"`
	ImdbID       string `json:"imdbID"`
	Year         string `json:"Year"`
	Genre        string `json:"Genre"`
	Language     string `json:"Language"`
	TotalSeasons string `json:"TotalSeasons"`
}

//SeasonsAPIContent definition
type SeasonsAPIContent struct {
	Season       string                `json:"Season"`
	TotalSeasons string                `json:"totalSeasons"`
	Episodes     []*EpisodesAPIContent `json:"Episodes"`
}

//EpisodesAPIContent definition
type EpisodesAPIContent struct {
	Title          string `json:"Title"`
	ImdbID         string `json:"ImdbID"`
	EpisodesNumber string `json:"Episode"`
	EpisodeContent MediaAPIContent
}

//Media definition
type Media struct {
	gorm.Model
	Type        MediaType  `json:"type"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Rating      string     `json:"rating"`
	Director    string     `json:"director"`
	Writer      string     `json:"writer"`
	Stars       string     `json:"stars"`
	ReleaseDate time.Time  `json:"releasedate"`
	Duration    string     `json:"duration"`
	ImdbID      string     `json:"imdbid"`
	Year        string     `json:"year"`
	Genre       string     `json:"genre"`
	Audio       string     `json:"audio"`
	Subtitles   string     `json:"subtitles"`
	Seasons     []*Seasons `gorm:"onDelete:CASCADE" json:"seasons"`
}

//Seasons definition
type Seasons struct {
	gorm.Model
	Season       int         `json:"season"`
	TotalSeasons int         `json:"totalSeasons"`
	Episode      []*Episodes `gorm:"onDelete:CASCADE" json:"episodes"`
	MediaID      *uint       `gorm:"not null" json:"mediaId"`
	Media        *Media
}

//Episodes definition
type Episodes struct {
	gorm.Model
	Episode   string `json:"episode"`
	Media     *Media `gorm:"onDelete:CASCADE" json:"content"`
	MediaID   *uint  `gorm:"not null" json:"mediaId"`
	SeasonsID *uint  `gorm:"not null" json:"seasonsId"`
	Seasons   *Seasons
}

//User definition
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//UserMedia definition
type UserMedia struct {
	gorm.Model
	MediaID *uint `gorm:"not null" json:"mediaId"`
	UserID  *uint `gorm:"not null" json:"userId"`
	Media   *Media
}

//Authentication definition
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Token definition
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

//New creates new service
func New(db *gorm.DB) Manager {
	db.AutoMigrate(&User{}, &Seasons{}, &Episodes{}, &Media{}, &UserMedia{})
	return &Data{DB: db}
}

//AddMovie adds movie to datastore
func (d *Data) AddMovie(body []byte) error {
	post := Media{}
	err := json.Unmarshal(body, &post)
	if err != nil {
		logger.Error.Println(err)
		return err
	}

	err = d.DB.Create(&post).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	d.DB.Save(&post)
	return err
}

//GetMovies gets all movies from datastore with given filters
func (d *Data) GetMovies(name, genre string) ([]Media, error) {
	result := []Media{}

	if name != "" {
		d.DB.Where("title = ?", name)
	}
	if genre != "" {
		d.DB.Where("genre LIKE ?", "%"+genre+"%")
	}
	err := d.DB.Where("type = ?", Movie).Find(&result).Limit(config.PageSize).Error

	return result, err
}

//GetMovieByID gets movie from datastore with given id
func (d *Data) GetMovieByID(id string) (Media, error) {
	result := Media{}
	err := d.DB.Where("type = ?", Movie).Where("id = ?", id).Find(&result).Error
	return result, err
}

//AddSeries adds series to datastore
func (d *Data) AddSeries(body []byte) error {
	post := Media{}
	err := json.Unmarshal(body, &post)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	err = d.DB.Create(&post).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	d.DB.Save(&post)
	return err
}

//DeleteMediaByID deletes series from datasource
func (d *Data) DeleteMediaByID(key string) error {
	err := d.DB.Transaction(func(tx *gorm.DB) error {
		result := Media{}
		id, err := strconv.Atoi(key)
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		err = d.DB.Preload("Seasons").Preload("Seasons.Episode").Preload("Seasons.Episode.Media").Where("id = ?", id).Find(&result).Error
		if err != nil {
			logger.Error.Println(err)
			return err
		}
		if result.Seasons != nil && len(result.Seasons) > 0 {
			for i := 0; i < len(result.Seasons); i++ {
				for j := 0; j < len(result.Seasons[i].Episode); j++ {
					err = d.DB.Delete(&Media{Model: gorm.Model{ID: *result.Seasons[i].Episode[j].MediaID}}).Error
					if err != nil {
						logger.Error.Println(err)
						return err
					}
				}
				err = d.DB.Where("seasons_id = ?", result.Seasons[i].ID).Delete(&Episodes{}).Error
				if err != nil {
					logger.Error.Println(err)
					return err
				}
			}
			err = d.DB.Where("media_id = ?", *result.Seasons[0].MediaID).Delete(&Seasons{}).Error
			if err != nil {
				logger.Error.Println(err)
				return err
			}
		}
		return d.DB.Delete(&Media{Model: gorm.Model{ID: uint(id)}}).Error
	})
	return err
}

//GetSeries gets all series from datastore with given filters
func (d *Data) GetSeries(name, genre string) ([]Media, error) {
	result := []Media{}

	if name != "" {
		d.DB.Where("title = ?", name)
	}
	if genre != "" {
		d.DB.Where("genre LIKE ?", "%"+genre+"%")
	}

	err := d.DB.Where("type = ?", Series).Find(&result).Limit(config.PageSize).Error
	return result, err
}

//GetSeriesByID gets series from datastore with given id
func (d *Data) GetSeriesByID(id string) (Media, error) {
	result := Media{}

	err := d.DB.Preload("Seasons").Preload("Seasons.Episode").Preload("Seasons.Episode.Media").Where("type = ?", Series).Where("id = ?", id).Find(&result).Error
	return result, err
}

//ConvertToMedia coverts response to madia
func (d *Data) ConvertToMedia(fromAPIContent MediaAPIContent, fromAPISeasons []SeasonsAPIContent) *Media {
	media := &Media{}
	if fromAPIContent.Type == "movie" {
		media.Type = Movie
	} else if fromAPIContent.Type == "series" {
		media.Type = Series
	} else if fromAPIContent.Type == "episode" {
		media.Type = Episode
	}
	media.Title = fromAPIContent.Title
	media.Description = fromAPIContent.Plot
	media.Genre = fromAPIContent.Genre
	media.Audio = fromAPIContent.Language
	media.Subtitles = fromAPIContent.Language
	media.Director = fromAPIContent.Director
	media.Writer = fromAPIContent.Writer
	media.Duration = fromAPIContent.Runtime
	media.ImdbID = fromAPIContent.ImdbID
	media.Rating = fromAPIContent.ImdbRating
	media.Year = fromAPIContent.Year
	media.Stars = fromAPIContent.Aktors
	if len(fromAPISeasons) > 0 {
		for i := 0; i < len(fromAPISeasons); i++ {
			seasons := &Seasons{}
			var err error
			seasons.Season, err = strconv.Atoi(fromAPISeasons[i].Season)
			if err != nil {
				seasons.Season = 0
			}
			seasons.TotalSeasons, err = strconv.Atoi(fromAPISeasons[i].Season)
			if err != nil {
				seasons.TotalSeasons = 0
			}
			if fromAPISeasons[i].Episodes != nil {
				for j := 0; j < len(fromAPISeasons[i].Episodes); j++ {
					episode := &Episodes{}
					episode.Episode = fromAPISeasons[i].Episodes[j].EpisodesNumber
					episode.Media = d.ConvertToMedia(fromAPISeasons[i].Episodes[j].EpisodeContent, nil)
					seasons.Episode = append(seasons.Episode, episode)
				}
			}
			media.Seasons = append(media.Seasons, seasons)
		}
	} else if media.Type == Series {
		seasons := &Seasons{}
		media.Seasons = []*Seasons{seasons}
	}

	return media
}

//ConvertToAPIContent converts to suggetsions API struct content
func (d *Data) ConvertToAPIContent(body []byte) (MediaAPIContent, error) {
	result := MediaAPIContent{}
	err := json.Unmarshal(body, &result)
	return result, err
}

//ConvertToAPISeasonsContent converts to suggetsions API struct content
func (d *Data) ConvertToAPISeasonsContent(body []byte) (SeasonsAPIContent, error) {
	result := SeasonsAPIContent{}
	err := json.Unmarshal(body, &result)
	return result, err
}

//AddFavorite adds series to datastore
func (d *Data) AddFavorite(body []byte) error {
	post := UserMedia{}
	err := json.Unmarshal(body, &post)
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	err = d.DB.Create(&post).Error
	if err != nil {
		logger.Error.Println(err)
		return err
	}
	d.DB.Save(&post)
	return err
}

//DeleteFavoriteByID adds series to datastore
func (d *Data) DeleteFavoriteByID(key string) error {
	id, err := strconv.Atoi(key)
	if err != nil {
		return err
	}
	err = d.DB.Delete(&UserMedia{Model: gorm.Model{ID: uint(id)}}).Error
	return err
}

//GetFavorites gets favorites medias from datastore with given user id
func (d *Data) GetFavorites(userID, name, genre string) ([]UserMedia, error) {
	result := []UserMedia{}

	if name != "" {
		d.DB.Where("title = ?", name)
	}
	if genre != "" {
		d.DB.Where("genre LIKE ?", "%"+genre+"%")
	}

	err := d.DB.Preload("Media").Where("user_id = ?", userID).Find(&result).Limit(config.PageSize).Error
	return result, err
}

//GetToken gets token for given valid user information
func (d *Data) GetToken(body []byte) (Token, error) {
	var authDetails Authentication
	err := json.Unmarshal(body, &authDetails)
	if err != nil {
		logger.Error.Println(err)
		return Token{}, err
	}

	var authUser User
	err = d.DB.Where("email =   ?", authDetails.Email).First(&authUser).Error

	if err != nil {
		logger.Error.Println(err)
		return Token{}, err
	}

	check := checkPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		return Token{}, errors.New(types.UsernamePasswordError)
	}

	validToken, err := generateJWT(authUser.Email, authUser.Role)
	if err != nil {
		logger.Error.Println(err)
		return Token{}, err
	}

	var token Token
	token.Email = authUser.Email
	token.Role = authUser.Role
	token.TokenString = validToken
	return token, err
}

func generateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(config.SecretKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//compare plain password with hash password
func checkPasswordHash(password, hash string) bool {
	return password == hash
}
