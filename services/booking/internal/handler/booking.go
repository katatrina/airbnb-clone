package handler

import (
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/pkg/middleware"
	"github.com/katatrina/airbnb-clone/pkg/request"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
)

const dateLayout = "2006-01-02"

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID

	var req CreateBookingRequest
	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	if _, err := uuid.Parse(req.ListingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"Invalid listing ID format")
		return
	}

	checkIn, err := time.Parse(dateLayout, req.CheckInDate)
	if err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"checkInDate must be in YYYY-MM-DD format")
		return
	}

	checkOut, err := time.Parse(dateLayout, req.CheckOutDate)
	if err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"checkOutDate must be in YYYY-MM-DD format")
		return
	}

	booking, err := h.bookingService.CreateBooking(c.Request.Context(),
		model.CreateBookingParams{
			ListingID:    req.ListingID,
			GuestID:      userID,
			CheckInDate:  checkIn,
			CheckOutDate: checkOut,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidDateRange):
			response.BadRequest(c, response.CodeValidationFailed,
				"Check-out date must be after check-in date")
		case errors.Is(err, model.ErrCheckInPast):
			response.BadRequest(c, response.CodeValidationFailed,
				"Check-in date cannot be in the past")
		case errors.Is(err, model.ErrListingNotFound):
			response.NotFound(c, response.CodeListingNotFound,
				"Listing not found or not available")
		case errors.Is(err, model.ErrSelfBooking):
			response.BadRequest(c, response.CodeValidationFailed,
				"You cannot book your own listing")
		case errors.Is(err, model.ErrDatesUnavailable):
			response.Conflict(c, response.CodeDatesUnavailable,
				"Selected dates are not available")
		case errors.Is(err, model.ErrListingServiceUnavailable):
			response.ServiceUnavailable(c,
				"Unable to verify listing. Please try again later")
		default:
			log.Printf("[ERROR] failed to create booking: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.Created(c, NewBookingResponse(booking),
		"Booking created successfully")
}

func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID
	bookingID := c.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"Invalid booking ID format")
		return
	}

	booking, err := h.bookingService.ConfirmBooking(
		c.Request.Context(), bookingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBookingNotFound),
			errors.Is(err, model.ErrNotBookingHost):
			response.NotFound(c, response.CodeBookingNotFound,
				"Booking not found")
		case errors.Is(err, model.ErrBookingNotPending):
			response.BadRequest(c, response.CodeBookingNotPending,
				"Only pending bookings can be confirmed")
		default:
			log.Printf("[ERROR] failed to confirm booking: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, NewBookingResponse(booking),
		"Booking confirmed successfully")
}

func (h *BookingHandler) RejectBooking(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID
	bookingID := c.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"Invalid booking ID format")
		return
	}

	booking, err := h.bookingService.RejectBooking(
		c.Request.Context(), bookingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBookingNotFound),
			errors.Is(err, model.ErrNotBookingHost):
			response.NotFound(c, response.CodeBookingNotFound,
				"Booking not found")
		case errors.Is(err, model.ErrBookingNotPending):
			response.BadRequest(c, response.CodeBookingNotPending,
				"Only pending bookings can be rejected")
		default:
			log.Printf("[ERROR] failed to reject booking: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, NewBookingResponse(booking),
		"Booking rejected successfully")
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID
	bookingID := c.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"Invalid booking ID format")
		return
	}

	booking, err := h.bookingService.CancelBooking(
		c.Request.Context(), bookingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrBookingNotFound),
			errors.Is(err, model.ErrNotBookingGuest):
			response.NotFound(c, response.CodeBookingNotFound,
				"Booking not found")
		case errors.Is(err, model.ErrBookingNotPending):
			response.BadRequest(c, response.CodeBookingNotPending,
				"Only pending bookings can be cancelled")
		default:
			log.Printf("[ERROR] failed to cancel booking: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.OK(c, NewBookingResponse(booking),
		"Booking cancelled successfully")
}

func (h *BookingHandler) GetBooking(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID
	bookingID := c.Param("id")

	if _, err := uuid.Parse(bookingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed,
			"Invalid booking ID format")
		return
	}

	booking, err := h.bookingService.GetBookingByID(
		c.Request.Context(), bookingID, userID)
	if err != nil {
		if errors.Is(err, model.ErrBookingNotFound) {
			response.NotFound(c, response.CodeBookingNotFound,
				"Booking not found")
			return
		}
		log.Printf("[ERROR] failed to get booking: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, NewBookingResponse(booking), "")
}

func (h *BookingHandler) ListGuestBookings(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID

	bookings, err := h.bookingService.ListGuestBookings(
		c.Request.Context(), userID)
	if err != nil {
		log.Printf("[ERROR] failed to list guest bookings: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, NewBookingsResponse(bookings), "")
}

func (h *BookingHandler) ListHostBookings(c *gin.Context) {
	userID := middleware.MustGetAuthUser(c).ID

	bookings, err := h.bookingService.ListHostBookings(
		c.Request.Context(), userID)
	if err != nil {
		log.Printf("[ERROR] failed to list host bookings: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, NewBookingsResponse(bookings), "")
}
