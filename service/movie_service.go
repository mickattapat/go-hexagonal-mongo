package service

import (
	"errors"
	"fmt"
	"golang-mongodb/repository"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type movieService struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return movieService{movieRepo: movieRepo}
}

func (s movieService) GetMovie(id int) (*MovieResponse, error) {
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

func (s movieService) GetMovies(id int) ([]MovieResponse, error) {
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

func (s movieService) GetAllMovies() ([]MovieResponse, error) {

	counter := sync.WaitGroup{}
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()
	movieSrv, err := s.movieRepo.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(len(movieSrv))
	ch := make(chan MovieResponse, len(movieSrv))
	responses := []MovieResponse{}
	for _, movie := range movieSrv {
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
	return responses, nil
}

func (s movieService) NewMovie(request NewMovieRequest) (*NewMovieResponse, error) {
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

func (s movieService) UpdateMovie(oid primitive.ObjectID, request NewMovieRequest) (*UpdateMovieResponse, error) {
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

func (s movieService) DeleteMovie(oid primitive.ObjectID) error {
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
