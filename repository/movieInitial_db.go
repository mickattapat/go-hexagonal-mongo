package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type movieRepositoryDB struct {
	db *mongo.Client
}

func NewMovieRepositoryDB(db *mongo.Client) MovieRepository {
	return movieRepositoryDB{db: db}
}

func (r movieRepositoryDB) GetByID(id int) (*Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movie := Movie{}
	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")

	filter := bson.M{"imdbID": id}
	err := collection.FindOne(ctx, filter).Decode(&movie)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &movie, nil
}

func (r movieRepositoryDB) GetByIDs(id int) ([]Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")

	movie := []Movie{}
	filter := bson.D{{Key: "imdbID", Value: id}}
	contactCursor, err := collection.Find(ctx, filter)
	defer contactCursor.Close(ctx)

	err = contactCursor.All(ctx, &movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r movieRepositoryDB) GetAll() ([]Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")

	movie := []Movie{}
	contactCursor, err := collection.Find(ctx, bson.M{})
	defer contactCursor.Close(ctx)

	err = contactCursor.All(ctx, &movie)
	if err != nil {
		fmt.Println("test")
		return nil, err
	}
	// fmt.Println(movie)

	return movie, nil
}

func (r movieRepositoryDB) Create(mv Movie) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")
	filter := Movie{ImdbID: mv.ImdbID, Title: mv.Title, Language: mv.Language}

	result, err := collection.InsertOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	mv.ID = oid

	return &mv, nil
}

func (r movieRepositoryDB) Update(oid primitive.ObjectID, mv Movie) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")

	// update := bson.D{{"$set", bson.D{
	// 	{"imdbID", 1234},
	// 	{"title", "Mick attapat"},
	// }}}
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "language", Value: mv.Language},
				{Key: "title", Value: mv.Title},
			},
		},
	}
	// movie := Movie{Title: mv.Title, Language: mv.Language}
	err := collection.FindOneAndUpdate(ctx, filter, update).Err()
	if err != nil {
		return nil, err
	}
	mv.ID = oid

	return &mv, nil
}

func (r movieRepositoryDB) Delete(oid primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := r.db.Database("uncleBobDvd").Collection("movie_initial")

	filter := bson.D{{Key: "_id", Value: oid}}
	result, err := collection.DeleteOne(ctx, &filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New(fmt.Sprintf("%v", result.DeletedCount))
	}
	return nil
}
