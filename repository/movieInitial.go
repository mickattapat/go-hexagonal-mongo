package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ImdbID      int                `json:"imdbID" bson:"imdbID"`
	Title       string             `json:"title"`
	Year        int                `json:"year" bson:"year"`
	Rating      string             `json:"rating"`
	Runtime     string             `json:"runtime"`
	Genre       string             `json:"genre"`
	Released    string             `json:"released"`
	Director    string             `json:"director"`
	Writer      string             `json:"writer"`
	Cast        string             `json:"cast"`
	Metacritic  string             `json:"metacritic"`
	ImdbRating  string             `json:"imdbRating"`
	ImdbVotes   int                `json:"imdbVotes" bson:"imdbVotes"`
	Poster      string             `json:"poster"`
	Plot        string             `json:"plot"`
	Fullplot    string             `json:"fullplot"`
	Language    string             `json:"language"`
	Country     string             `json:"country"`
	Awards      string             `json:"awards"`
	Lastupdated string             `json:"lastupdated"`
	Type        string             `json:"type"`
}

type MovieRepository interface {
	GetByID(int) (*Movie, error) //ถ้าใช้ struct (Customer) จะ return nill กลับไปไม่ได้จึงต้องใช้ pointer
	GetByIDs(int) ([]Movie, error)
	GetAll() ([]Movie, error)
	Create(Movie) (*Movie, error)
	Update(primitive.ObjectID, Movie) (*Movie, error)
	Delete(primitive.ObjectID) error
}
