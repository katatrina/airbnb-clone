package handler

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/katatrina/airbnb-clone/pkg/request"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/listing/internal/constant"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (h *ListingHandler) ListActiveListings(c *gin.Context) {
	paginationParams := request.ParsePaginationParams(c)

	// TODO: Add filtering and searching

	listings, total, err := h.listingService.ListActiveListings(
		c.Request.Context(),
		paginationParams.Limit(),
		paginationParams.Offset(),
	)
	if err != nil {
		log.Printf("[ERROR] failed to list active listings: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OKWithPagination(c, NewListingsResponse(listings), "", paginationParams.Page, paginationParams.PageSize, total)
}

func (h *ListingHandler) GetActiveListing(c *gin.Context) {
	listingID := c.Param("id")

	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid listing ID format")
		return
	}

	listing, err := h.listingService.GetActiveListingByID(c.Request.Context(), listingID)
	if err != nil {
		if errors.Is(err, model.ErrListingNotFound) {
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
			return
		}

		log.Printf("[ERROR] failed to get listing by ID: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, NewListingResponse(listing), "")
}

func (h *ListingHandler) CreateListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)

	var req CreateListingRequest

	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	listing, err := h.listingService.CreateListing(c.Request.Context(), model.CreateListingParams{
		HostID:        userID,
		Title:         req.Title,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
		ProvinceCode:  req.ProvinceCode,
		DistrictCode:  req.DistrictCode,
		WardCode:      req.WardCode,
		AddressDetail: req.AddressDetail,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrProvinceCodeNotFound):
			response.BadRequest(c, response.CodeProvinceNotFound, "Province code not found")
		case errors.Is(err, model.ErrDistrictCodeNotFound):
			response.BadRequest(c, response.CodeDistrictNotFound, "District code not found")
		case errors.Is(err, model.ErrWardCodeNotFound):
			response.BadRequest(c, response.CodeWardNotFound, "Ward code not found")
		case errors.Is(err, model.ErrDistrictProvinceMismatch):
			response.BadRequest(c, response.CodeReferenceInvalid, "District does not belong to province")
		case errors.Is(err, model.ErrWardDistrictMismatch):
			response.BadRequest(c, response.CodeReferenceInvalid, "Ward does not belong to district")
		default:
			log.Printf("[ERROR] failed to create listing: %v", err)
			response.InternalServerError(c)
		}
		return
	}

	response.Created(c, NewListingResponse(listing), "Listing created successfully")
}

func (h *ListingHandler) UpdateListingBasicInfo(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")
	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid listing ID format")
		return
	}

	var req UpdateListingBasicInfoRequest
	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	listing, err := h.listingService.UpdateListingBasicInfo(c.Request.Context(), listingID, userID, model.UpdateListingBasicInfoParams{
		Title:         req.Title,
		Description:   req.Description,
		PricePerNight: req.PricePerNight,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		case errors.Is(err, model.ErrActiveListingCannotBeUpdated):
			response.BadRequest(c, response.CodeActiveListingCannotBeUpdated, "Active listing cannot be updated")
		default:
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, NewListingResponse(listing), "Listing basic info updated successfully")
}

func (h *ListingHandler) UpdateListingAddress(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")
	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid listing ID format")
		return
	}

	var req UpdateListingAddressRequest
	if err := request.ShouldBindJSON(c, &req); err != nil {
		response.HandleJSONBindingError(c, err)
		return
	}

	listing, err := h.listingService.UpdateListingAddress(c.Request.Context(), model.UpdateListingAddressParams{
		ListingID:     listingID,
		HostID:        userID,
		ProvinceCode:  req.ProvinceCode,
		DistrictCode:  req.DistrictCode,
		WardCode:      req.WardCode,
		AddressDetail: req.AddressDetail,
	})
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		case errors.Is(err, model.ErrActiveListingCannotBeUpdated):
			response.BadRequest(c, response.CodeActiveListingCannotBeUpdated, "Active listing cannot be updated")
		case errors.Is(err, model.ErrProvinceCodeNotFound):
			response.BadRequest(c, response.CodeProvinceNotFound, "Province code not found")
		case errors.Is(err, model.ErrDistrictCodeNotFound):
			response.BadRequest(c, response.CodeDistrictNotFound, "District code not found")
		case errors.Is(err, model.ErrWardCodeNotFound):
			response.BadRequest(c, response.CodeWardNotFound, "Ward code not found")
		case errors.Is(err, model.ErrDistrictProvinceMismatch):
			response.BadRequest(c, response.CodeReferenceInvalid, "District does not belong to province")
		case errors.Is(err, model.ErrWardDistrictMismatch):
			response.BadRequest(c, response.CodeReferenceInvalid, "Ward does not belong to district")
		default:
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, NewListingResponse(listing), "Listing address updated successfully")
}

func (h *ListingHandler) PublishListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")
	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid listing ID format")
		return
	}

	listing, err := h.listingService.PublishListing(c.Request.Context(), listingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		case errors.Is(err, model.ErrListingNotDraft):
			response.BadRequest(c, response.CodeListingNotDraft, "Listing must be in draft status to publish")
		case errors.Is(err, model.ErrListingIncomplete):
			response.BadRequest(c, response.CodeListingIncomplete, "Listing is incomplete. Please add full required information before publishing (min 50 characters)")
		default:
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, NewListingResponse(listing), "Listing published successfully")
}

func (h *ListingHandler) DeactivateListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")

	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Listing ID format invalid")
		return
	}

	listing, err := h.listingService.DeactivateListingByID(c.Request.Context(), listingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		case errors.Is(err, model.ErrListingNotActive):
			response.BadRequest(c, response.CodeListingNotActive, "Listing must be in active status to deactivate")
		default:
			log.Printf("[ERROR] failed to deactivate host listing: %v", err)
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, listing, "Deactivate listing successfully")
}

func (h *ListingHandler) ReactivateListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")

	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Listing ID format invalid")
		return
	}

	listing, err := h.listingService.ReactivateListingByID(c.Request.Context(), listingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		case errors.Is(err, model.ErrListingNotInactive):
			response.BadRequest(c, response.CodeListingNotInactive, "Listing must be in inactive status to reactivate")
		default:
			log.Printf("[ERROR] failed to reactivate host listing: %v", err)
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, listing, "Reactivate listing successfully")
}

func (h *ListingHandler) DeleteListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")

	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Listing ID format invalid")
		return
	}

	err := h.listingService.DeleteListingByID(c.Request.Context(), listingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		default:
			log.Printf("[ERROR] failed to delete host listing: %v", err)
			response.InternalServerError(c)
		}

		return
	}

	response.NoContent(c)
}

func (h *ListingHandler) ListHostListings(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)

	// TODO: Add filtering and searching
	listings, err := h.listingService.ListHostListings(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[ERROR] failed to list host listings: %v", err)
		response.InternalServerError(c)
		return
	}

	response.OK(c, NewListingsResponse(listings), "")
}

func (h *ListingHandler) GetUserListing(c *gin.Context) {
	userID := c.MustGet(constant.UserIDKey).(string)
	listingID := c.Param("id")

	if _, err := uuid.Parse(listingID); err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Listing ID format invalid")
		return
	}

	listing, err := h.listingService.GetHostListingByID(c.Request.Context(), listingID, userID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrListingNotFound), errors.Is(err, model.ErrListingOwnerMismatch):
			response.NotFound(c, response.CodeListingNotFound, "Listing not found")
		default:
			log.Printf("[ERROR] failed to get host listing: %v", err)
			response.InternalServerError(c)
		}

		return
	}

	response.OK(c, NewListingResponse(listing), "")
}
