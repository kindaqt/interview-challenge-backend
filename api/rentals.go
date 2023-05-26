package api

import (
	"fmt"
	"log"
	"net/http"
	"outdoorsy/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Price struct {
	Day int `json:"day"`
}

type Location struct {
	City    string  `json:"city"`
	State   string  `json:"state"`
	Zip     string  `json:"zip"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Rental struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Type            string   `json:"type"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Year            int      `json:"year"`
	Length          float64  `json:"length"`
	Sleeps          int      `json:"sleeps"`
	PrimaryImageURL string   `json:"primary_image_url"`
	Price           Price    `json:"price"`
	Location        Location `json:"location"`
	User            User     `json:"user"`
}

func (r *Router) getRentalByID(c *gin.Context) {
	// Parse rentalID
	rentalIDParam := c.Param("rentalID")
	if rentalIDParam == "" {
		err := fmt.Errorf(`invalid rentalID query param cannot be blank: rentalID="%s"`, rentalIDParam)
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	rentalID, err := strconv.Atoi(rentalIDParam)
	if err != nil {
		err := fmt.Errorf(`invalid rentalID query param must be an int: rentalID="%s"`, rentalIDParam)
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	// Connect to DB
	conn, err := r.db.Connect()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database",
		})
	}
	defer conn.Close()

	// Start a transaction
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start database transaction",
		})
	}

	// Query rental by id
	rental, err := r.db.QueryRentalByID(tx, rentalID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rental",
		})
	}

	// Query user by id
	user, err := r.db.QueryUserByID(tx, rental.UserID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
	}

	body := Rental{
		ID:              rental.ID,
		Name:            rental.Name,
		Description:     rental.Description,
		Type:            rental.Type,
		Make:            rental.VehicleMake,
		Model:           rental.VehicleModel,
		Year:            rental.VehicleYear,
		Length:          rental.VehicleLength,
		Sleeps:          rental.Sleeps,
		PrimaryImageURL: rental.PrimaryImageURL,
		Price: Price{
			Day: int(rental.PricePerDay),
		},
		Location: Location{
			City:    rental.HomeCity,
			State:   rental.HomeState,
			Zip:     rental.HomeZip,
			Country: rental.HomeCountry,
			Lat:     rental.Lat,
			Lng:     rental.Lng,
		},
		User: User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	c.IndentedJSON(http.StatusOK, body)
}

type QueryRentalsParams struct {
	PriceMin *int `form:"price_min" binding:"omitempty,min=0"`
	PriceMax *int `form:"price_max" binding:"omitempty,min=0"`
	Limit    *int `form:"limit" binding:"omitempty,min=1"`
	Offset   *int `form:"offset" binding:"omitempty,min=0"`
	// TODO: validate
	IDs string `form:"ids"`
	// TODO: validate
	Near string `form:"near"`
	// TODO: validate
	Sort string `form:"sort"`
}

func (r *Router) getRentals(c *gin.Context) {
	var params QueryRentalsParams

	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	queryParams := db.QueryRentalsParams{
		PriceMin: params.PriceMin,
		PriceMax: params.PriceMax,
		Limit:    params.Limit,
		Offset:   params.Offset,
		IDs:      params.IDs,
		Near:     params.Near,
		Sort:     params.Sort,
	}

	// TODO: add validation

	// Connect to DB
	conn, err := r.db.Connect()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to database",
		})
	}
	defer conn.Close()

	// Start a transaction
	tx, err := conn.Begin()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start database transaction",
		})
	}

	// Query rentals
	rentals, err := r.db.QueryRentals(tx, &queryParams)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch rentals",
		})
	}

	var body []Rental
	for _, rental := range rentals {
		// Query user by id
		user, err := r.db.QueryUserByID(tx, rental.UserID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch user",
			})
		}

		body = append(body, Rental{
			ID:              rental.ID,
			Name:            rental.Name,
			Description:     rental.Description,
			Type:            rental.Type,
			Make:            rental.VehicleMake,
			Model:           rental.VehicleModel,
			Year:            rental.VehicleYear,
			Length:          rental.VehicleLength,
			Sleeps:          rental.Sleeps,
			PrimaryImageURL: rental.PrimaryImageURL,
			Price: Price{
				Day: int(rental.PricePerDay),
			},
			Location: Location{
				City:    rental.HomeCity,
				State:   rental.HomeState,
				Zip:     rental.HomeZip,
				Country: rental.HomeCountry,
				Lat:     rental.Lat,
				Lng:     rental.Lng,
			},
			User: User{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
		})
	}

	c.IndentedJSON(http.StatusOK, body)
}
