package main

import (
	"context"
	"fmt"
	"golang-mongodb/handler"
	"golang-mongodb/repository"
	"golang-mongodb/service"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	initTimeZone()
	initConfig()
	db := initDB()
	app := fiber.New()
	redisClient := initRedis()
	_ = redisClient
	// repo
	movieRepository := repository.NewMovieRepositoryDB(db)
	// service
	// movieService := service.NewMovieService(movieRepository)
	movieService := service.NewMovieServiceRedis(movieRepository, redisClient)
	// handler
	movieHandler := handler.NewMovieHandler(movieService)

	// logger
	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Bangkok",
	}))

	// api
	app.Get("/movie/:movieId", movieHandler.GetMovie)
	app.Get("/movies/:movieId", movieHandler.GetMovies)
	app.Get("/movies", movieHandler.GeAlltMovies)
	app.Post("/movie", movieHandler.CreateMovie)
	app.Put("/movie/:oid", movieHandler.UpdateMovie)
	app.Delete("/movie/:oid", movieHandler.DeleteMovie)

	// listen port
	appPort := fmt.Sprintf(":%v", viper.GetInt("app.port"))
	app.Listen(appPort)
}

func initTimeZone() {
	// set timezone ให้กับระบบ ป้องกันการมีปัญหาหากไปใช้ container
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func initDB() *mongo.Client {
	dns := fmt.Sprintf("%v://%v:%v", viper.GetString("db.driver"), viper.GetString("db.host"), viper.GetInt("db.port"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(dns))
	if err != nil {
		log.Fatal("cannot connect to DB !", err)
	}

	err = db.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("cannot connect to DB !", err)
	}
	return db
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	// SET ENV
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// Ex APP_PORT=3000 go run .

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
