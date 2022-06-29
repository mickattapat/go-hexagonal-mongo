package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-mongodb/repository"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type movieServiceRedis struct {
	movieRepo   repository.MovieRepository
	redisClient *redis.Client
}

func NewMovieServiceRedis(movieRepo repository.MovieRepository, redisClient *redis.Client) MovieService {
	return movieServiceRedis{movieRepo, redisClient}
}

func (s movieServiceRedis) GetAllMovies() ([]MovieResponse, error) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	key := "service::GetUsers"

	mv := []MovieResponse{}
	// Redis Get
	if productJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(productJson), &mv) == nil {
			fmt.Println("redis")
			return mv, nil
		}
	}

	movieSrv, err := s.movieRepo.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(len(movieSrv))

	responses := []MovieResponse{}
	for _, movie := range movieSrv {
		movie := movie
		responses = append(responses, MovieResponse{
			ID:         movie.ID,
			Title:      movie.Title,
			ImdbID:     movie.ImdbID,
			Year:       movie.Year,
			Rating:     movie.Rating,
			Released:   movie.Released,
			Director:   movie.Director,
			Writer:     movie.Writer,
			Language:   movie.Language,
			Awards:     movie.Awards,
			ImdbVotes:  movie.ImdbVotes,
			ImdbRating: movie.Rating,
		})
	}

	if len(responses) <= 0 {
		return nil, errors.New("movie not found !")
	}
	fmt.Println("2", len(responses))

	// Redis Set
	if data, err := json.Marshal(responses); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*60)
	}
	fmt.Println("database")

	return responses, nil
}

func (s movieServiceRedis) GetMovie(id int) (*MovieResponse, error) {
	movieSrv, err := s.movieRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("movie not found !")
		}

		return nil, errors.New(err.Error())
	}
	movieResponse := MovieResponse{
		ID:         movieSrv.ID,
		Title:      movieSrv.Title,
		ImdbID:     movieSrv.ImdbID,
		Year:       movieSrv.Year,
		Rating:     movieSrv.Rating,
		Released:   movieSrv.Released,
		Director:   movieSrv.Director,
		Writer:     movieSrv.Writer,
		Language:   movieSrv.Language,
		Awards:     movieSrv.Awards,
		ImdbVotes:  movieSrv.ImdbVotes,
		ImdbRating: movieSrv.ImdbRating,
	}

	return &movieResponse, nil
}

func (s movieServiceRedis) GetMovies(id int) ([]MovieResponse, error) {
	movieSrv, err := s.movieRepo.GetByIDs(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("movie not found !")
		}

		return nil, errors.New(err.Error())
	}
	responses := []MovieResponse{}
	for _, movie := range movieSrv {
		responses = append(responses, MovieResponse{
			ID:         movie.ID,
			Title:      movie.Title,
			ImdbID:     movie.ImdbID,
			Year:       movie.Year,
			Rating:     movie.Rating,
			Released:   movie.Released,
			Director:   movie.Director,
			Writer:     movie.Writer,
			Language:   movie.Language,
			Awards:     movie.Awards,
			ImdbVotes:  movie.ImdbVotes,
			ImdbRating: movie.ImdbRating,
		})
	}
	if len(responses) <= 0 {
		return nil, errors.New("movie not found !")
	}

	return responses, nil
}

func (s movieServiceRedis) NewMovie(request NewMovieRequest) (*NewMovieResponse, error) {
	// validate
	// data
	movie := repository.Movie{
		Title:    request.Title,
		Language: request.Language,
	}

	movieNew, err := s.movieRepo.Create(movie)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	response := NewMovieResponse{
		ImdbID:   movieNew.ImdbID,
		Title:    movieNew.Title,
		Language: movieNew.Language,
	}

	return &response, nil
}

func (s movieServiceRedis) UpdateMovie(oid primitive.ObjectID, request NewMovieRequest) (*UpdateMovieResponse, error) {
	movie := repository.Movie{
		Title:    request.Title,
		Language: request.Language,
	}

	movieUpdate, err := s.movieRepo.Update(oid, movie)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("movie not found !")
		}

		return nil, errors.New(err.Error())
	}

	response := UpdateMovieResponse{
		ID:       movieUpdate.ID,
		Title:    movieUpdate.Title,
		Language: movieUpdate.Language,
	}
	return &response, nil
}
func (s movieServiceRedis) DeleteMovie(oid primitive.ObjectID) error {
	err := s.movieRepo.Delete(oid)
	errnew := fmt.Sprintf("%v", err)
	if errnew == "0" {
		return errors.New("id is not found !")
	}
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
