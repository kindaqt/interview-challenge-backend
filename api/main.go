package api

import (
	"net/http"
	"outdoorsy/db"

	"github.com/gin-gonic/gin"
)

type Router struct {
	db *db.Database
}

func NewRouter(database *db.Database) (*gin.Engine, error) {
	router := gin.Default()

	api := &Router{
		db: database,
	}

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.GET("/rentals/:rentalID", api.getRentalByID)
	router.GET("/rentals", api.getRentals)

	return router, nil
}
