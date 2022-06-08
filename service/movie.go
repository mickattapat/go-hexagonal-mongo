package service

import "go.mongodb.org/mongo-driver/bson/primitive"

type MovieResponse struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      string             `json:"title"`
	ImdbID     int                `json:"imdbID" bson:"imdbID"`
	Year       int                `json:"year" bson:"year"`
	Rating     string             `json:"rating"`
	Released   string             `json:"released"`
	Director   string             `json:"director"`
	Writer     string             `json:"writer"`
	Language   string             `json:"language"`
	Awards     string             `json:"awards"`
	ImdbVotes  int                `json:"imdbVotes" bson:"imdbVotes"`
	ImdbRating string             `json:"imdbRating"`
}

type NewMovieRequest struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

type NewMovieResponse struct {
	// ID       primitive.ObjectID `json:"_id" bson:"_id"`
	ImdbID   int    `json:"imdbID" bson:"imdbID"`
	Title    string `json:"title"`
	Language string `json:"language"`
}
type UpdateMovieResponse struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Title    string             `json:"title"`
	Language string             `json:"language"`
}

type MovieService interface {
	GetMovie(int) (*MovieResponse, error)
	GetMovies(int) ([]MovieResponse, error)
	GetAllMovies() ([]MovieResponse, error)
	NewMovie(NewMovieRequest) (*NewMovieResponse, error)
	UpdateMovie(primitive.ObjectID, NewMovieRequest) (*UpdateMovieResponse, error)
	DeleteMovie(primitive.ObjectID) error
}
