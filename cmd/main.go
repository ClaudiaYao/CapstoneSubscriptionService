package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"time"

	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/data"
	domain "github.com/ClaudiaYao/CapstoneSubscriptionService/app/domain"
	"github.com/ClaudiaYao/CapstoneSubscriptionService/app/domain/auth"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

var counts int64

func main() {
	log.Println("Starting subscription service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	appCon := getConfig()
	tokenPath := "app/domain/auth/tokenData"
	jwtMaker, err := auth.NewJWTMaker(tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	jwtVerifier, err := auth.NewJWTVerifier(tokenPath)
	if err != nil {
		log.Fatal(err)
	}

	//set up
	subService := &domain.SubscriptionService{
		DBConnection: conn,
		AppConfig:    appCon,
		JwtMaker:     jwtMaker,
		JwtVerifier:  jwtVerifier,
	}

	srv := &http.Server{
		Addr:    appCon.ServicePort,
		Handler: subService.Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func getConfig() *domain.AppConfiguration {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	expireSec, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE_SECS"))
	if err != nil {
		log.Fatal(err)
	}

	return &domain.AppConfiguration{
		TokenExpireSecs:                  expireSec,
		ServicePort:                      os.Getenv("SERVICE_PORT"),
		EmailServiceContainerName:        os.Getenv("MAIL_SERVICE"),
		LoginServiceContainerName:        os.Getenv("LOGIN_SERVICE"),
		PlaylistServiceContainerName:     os.Getenv("PLATLIST_SERVICE"),
		SubscriptionServiceContainerName: os.Getenv("SUBSCRIPTION_SERVICE"),
	}
}

// C: this function will connect to database and then return *DataQuery
func openDB(dsn string) (*data.DataQuery, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dataquery := data.DataQuery{DBConn: db}

	return &dataquery, nil
}

// C: wrap the openDB function and provide retry mechanism
func connectToDB() *data.DataQuery {
	dsn := os.Getenv("DSN")

	println("debugging line:", dsn)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	// load .env file which is located at the root path
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
