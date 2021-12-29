package specs

import (
	"scaleflixapi/data"
)

//CreateTestMovie for tests
func CreateTestMovie() data.Media {
	return data.Media{
		Title:       "test movie",
		Description: "test movie desc",
		Type:        data.Movie,
		Year:        "2021",
		Director:    "test director",
		Writer:      "test writer",
		Stars:       "test stars",
		Genre:       "Action, Adventure, Drama",
	}
}

//CreateTestSeries for tests
func CreateTestSeries() data.Media {
	return data.Media{
		Title:       "test Series",
		Description: "test Series desc",
		Type:        data.Series,
		Year:        "2021",
		Director:    "test director",
		Writer:      "test writer",
		Stars:       "test stars",
		Genre:       "Action, Adventure, Drama",
		ImdbID:      "tt0944947",
		Seasons: []*data.Seasons{
			{
				Season:       1,
				TotalSeasons: 1,
				Episode: []*data.Episodes{
					{
						Episode: "1",
						Media: &data.Media{
							Title:       "test Episode",
							Description: "test Episode desc",
							Type:        data.Episode,
							Year:        "2021",
							Director:    "test director",
							Writer:      "test writer",
							Stars:       "test stars",
							Genre:       "Action, Adventure, Drama",
							ImdbID:      "tt1480055",
						}},
					{
						Episode: "2",
						Media: &data.Media{
							Title:       "test Episode 2",
							Description: "test Episode desc 2",
							Type:        data.Episode,
							Year:        "2021",
							Director:    "test director",
							Writer:      "test writer",
							Stars:       "test stars",
							Genre:       "Action, Adventure, Drama",
							ImdbID:      "tt1480056",
						},
					},
				},
			},
		},
	}
}

//CreateFavorites for test
func CreateFavorites() data.UserMedia {
	var s uint = 1
	return data.UserMedia{
		UserID:  &s,
		MediaID: &s,
	}
}

//CreateAdminUser creates admin user for test
func CreateAdminUser() data.User {
	return data.User{
		Name:     "user3",
		Email:    "user3@gmail.com",
		Password: "user3222",
		Role:     "admin",
	}
}

//CreateUser creates user for test
func CreateUser() data.User {
	return data.User{
		Name:     "user2",
		Email:    "user2@gmail.com",
		Password: "user2111",
		Role:     "user",
	}
}
