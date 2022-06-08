package handler

import (
	"golang-mongodb/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type movieHandler struct {
	movieHand service.MovieService
}

func NewMovieHandler(movieHand service.MovieService) movieHandler {
	return movieHandler{movieHand: movieHand}
}

func (h movieHandler) GetMovie(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("movieId")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}
	movie, err := h.movieHand.GetMovie(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.JSON(movie)
}

func (h movieHandler) GetMovies(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("movieId")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	movie, err := h.movieHand.GetMovies(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.JSON(movie)
}

func (h movieHandler) GeAlltMovies(ctx *fiber.Ctx) error {

	movie, err := h.movieHand.GetAllMovies()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.JSON(movie)
}

func (h movieHandler) CreateMovie(ctx *fiber.Ctx) error {

	request := service.NewMovieRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return err
	}

	if request.Language == "" || request.Title == "" {
		return fiber.ErrUnprocessableEntity
	}

	responses, err := h.movieHand.NewMovie(request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.JSON(responses)
}

func (h movieHandler) UpdateMovie(ctx *fiber.Ctx) error {
	id := ctx.Params("oid")

	request := service.NewMovieRequest{}
	err := ctx.BodyParser(&request)
	if err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	if request.Language == "" || request.Title == "" {
		return fiber.ErrUnprocessableEntity
	}

	responses, err := h.movieHand.UpdateMovie(oid, request)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.JSON(responses)
}

func (h movieHandler) DeleteMovie(ctx *fiber.Ctx) error {
	id := ctx.Params("oid")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	err = h.movieHand.DeleteMovie(oid)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": false, "error": err.Error()})
	}

	return ctx.Status(200).JSON(fiber.Map{"msg": true})
}
