package util

import (
	"fmt"
	"strings"
	"time"

	"tonx_task/service"
)

func Booking(maxRetries int, flightID, userID uint, fareClass string, price float64) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = service.FlightBookingService.CreateBooking(flightID, userID, fareClass, price)
		if err == nil {
			// Booking succeeded
			return nil
		}

		// Only retry if the error message indicates an optimistic lock conflict
		if !strings.Contains(err.Error(), "optimistic lock conflict") {
			return err
		}

		// Wait 0.5 sec, before next retrying
		time.Sleep(50 * time.Millisecond)
	}

	return fmt.Errorf("booking failed after %d retries: %v", maxRetries, err)
}
