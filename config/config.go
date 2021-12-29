package config

import (
	"scaleflixapi/utils"
)

var (
	//ExportFilePath definition
	ExportFilePath = utils.GetEnv("EXPORT_FILE_PATH", "./output")
	//APIPort definition
	APIPort = utils.GetEnv("API_PORT", ":8080")
	//DBHost definition
	DBHost = utils.GetEnv("DB_HOST", "localhost")
	//DBPort definition
	DBPort = utils.GetEnv("DB_PORT", "5454")
	//DBUser definition
	DBUser = utils.GetEnv("DB_USER", "postgres")
	//DBPassword definition
	DBPassword = utils.GetEnv("DB_PASSWORD", "postgres")
	//DBName definition
	DBName = utils.GetEnv("DB_DBNAME", "postgres")
	//PageSize definition
	PageSize = utils.GetEnv("PAGE_SIZE", "10")
	//SecretKey definition
	SecretKey = utils.GetEnv("SECRET_KEY", "secretkeyjwt")
	//DBNameTest definition
	DBNameTest = utils.GetEnv("DB_DBNAME_TEST", "postgrestest")
	//APIKey definition
	APIKey = utils.GetEnv("API_KEY", "*****")
)
