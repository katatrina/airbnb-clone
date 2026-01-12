package handler

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/pkg/validator"
	"github.com/katatrina/airbnb-clone/services/listing/internal/constant"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
	"github.com/katatrina/airbnb-clone/services/listing/internal/service"
)

func (h *ListingHandler) ListActiveListings(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) GetListingByID(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) CreateListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)

	var req CreateListingRequest

	if err := validator.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	listing, err := h.listingService.CreateListing(c.Request.Context(), service.CreateListingParams{
		HostID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
		ProvinceCode:  req.ProvinceCode,
		WardCode:      req.WardCode,
		AddressDetail: req.AddressDetail,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrProvinceCodeNotFound):
			response.BadRequest(c, response.CodeResourceNotFound, fmt.Sprintf("Province code %s not found", req.ProvinceCode))
		case errors.Is(err, model.ErrWardCodeNotFound):
			response.BadRequest(c, response.CodeResourceNotFound, fmt.Sprintf("Ward code %s not found", req.WardCode))
		case errors.Is(err, model.ErrWardProvinceMismatch):
			response.BadRequest(c, response.CodeReferenceInvalid, fmt.Sprintf("Ward with code %s does not belong to province with code %s", req.WardCode, req.ProvinceCode))
		default:
			log.Printf("[ERROR] failed to create listing: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.Created(c, NewListingResponse(listing))
}

func (h *ListingHandler) UpdateListing(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) PublishListing(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) DeactivateListing(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) DeleteListing(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) ListUserListings(c *gin.Context) {
	panic("not implemented")
}

func (h *ListingHandler) GetUserListingByID(c *gin.Context) {
	panic("not implemented")
}
