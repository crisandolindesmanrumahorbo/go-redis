package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-gin-redis/cache"
)

var (
	redisCache = cache.NewRedisCache("localhost:6379", 0, 1)
)

func main() {
	r := gin.Default()

	r.POST("/movies", createHandler)
	r.GET("/movies/:id", getMovieByIdHandler)
	r.GET("/movies", getMoviesHandler)
	r.PUT("/movies/:id", updateMovieHandler)
	r.DELETE("/movies/:id", deleteMovieHandler)

	r.Run(":5000")
}

func createHandler(ctx *gin.Context) {
	var movie cache.Movie
	err := ctx.ShouldBind(&movie)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res, err := redisCache.CreateMovie(movie)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"movie": res,
	})
	return
}

func getMovieByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := redisCache.GetMovie(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"movie": res,
	})
}

func getMoviesHandler(ctx *gin.Context) {
	res, err := redisCache.GetMovies()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"movies": res,
	})
}

func updateMovieHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := redisCache.GetMovie(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	var movie cache.Movie
	err = ctx.ShouldBind(&movie)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	res.Title = movie.Title
	res.Description = movie.Description
	result, err := redisCache.UpdateMovie(*res)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"movie": result,
	})
}

func deleteMovieHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := redisCache.DeleteMovie(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("movie with id %s deleted", id),
	})
}
