package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-mongodb/repository"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	counter := sync.WaitGroup{}
	movieSrv, err := s.movieRepo.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(len(movieSrv))
	ch := make(chan MovieResponse, len(movieSrv[:1000]))
	responses := []MovieResponse{}
	for _, movie := range movieSrv[:1000] {
		movie := movie
		counter.Add(1)
		go func() {
			defer counter.Done()
			ch <- MovieResponse{
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
			}
			// responses = append(responses, MovieResponse{
			// 	ID:         movie.ID,
			// 	Title:      movie.Title,
			// 	ImdbID:     movie.ImdbID,
			// 	Year:       movie.Year,
			// 	Rating:     movie.Rating,
			// 	Released:   movie.Released,
			// 	Director:   movie.Director,
			// 	Writer:     movie.Writer,
			// 	Language:   movie.Language,
			// 	Awards:     movie.Awards,
			// 	ImdbVotes:  movie.ImdbVotes,
			// 	ImdbRating: movie.Rating,
			// })
		}()
	}
	go func() {
		counter.Wait()
		close(ch)
	}()

	for v := range ch {
		responses = append(responses, MovieResponse{
			ID:         v.ID,
			Title:      v.Title,
			ImdbID:     v.ImdbID,
			Year:       v.Year,
			Rating:     v.Rating,
			Released:   v.Released,
			Director:   v.Director,
			Writer:     v.Writer,
			Language:   v.Language,
			Awards:     v.Awards,
			ImdbVotes:  v.ImdbVotes,
			ImdbRating: v.Rating,
		})
	}

	if len(responses) <= 0 {
		return nil, errors.New("movie not found !")
	}
	fmt.Println("2", len(responses))
	counter.Wait()

	// Redis Set
	if data, err := json.Marshal(responses); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*60)
	}
	fmt.Println("database")

	return responses, nil
}

func (s movieServiceRedis) GetMovie(id int) (*MovieResponse, error) {
	return nil, nil
}

func (s movieServiceRedis) GetMovies(id int) ([]MovieResponse, error) {
	return nil, nil
}

func (s movieServiceRedis) NewMovie(request NewMovieRequest) (*NewMovieResponse, error) {
	return nil, nil
}

func (s movieServiceRedis) UpdateMovie(oid primitive.ObjectID, request NewMovieRequest) (*UpdateMovieResponse, error) {
	return nil, nil
}

func (s movieServiceRedis) DeleteMovie(oid primitive.ObjectID) error {
	return nil
}
