package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"tonx_task/service"
	"tonx_task/util"

	"github.com/gin-gonic/gin"
)

func init() {
	// Obtain the root router group.
	root := GetRoot()

	// Create router group for user module.
	handlerGroup := root.Group("")
	handlerGroup.GET("flights", GetFlights)
	// TODO: middleware.AuthLogin() (user authentication)
	handlerGroup.POST("bookings", CreateBooking)

}

func GetFlights(ctx *gin.Context) {
	origin := ctx.Query("origin")
	destination := ctx.Query("destination")
	dateStr := ctx.Query("date")

	// default page = 1, page_size = 10
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}

	var flightDate time.Time
	if dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			flightDate = t
		}
	}

	flights, err := service.FlightBookingService.SearchFlights(origin, destination, flightDate, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "result": flights})
}

func CreateBooking(ctx *gin.Context) {
	var req struct {
		FlightID  uint    `json:"flight_id"`
		UserID    uint    `json:"user_id"`
		FareClass string  `json:"fare_class"`
		Price     float64 `json:"price"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("invalid request body: %v\n", err)
		return
	}

	maxRetries := 3

	if err := util.Booking(maxRetries, req.FlightID, req.UserID, req.FareClass, req.Price); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "booking created"})
}
