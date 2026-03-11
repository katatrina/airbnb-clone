package handler

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
	"github.com/katatrina/airbnb-clone/services/listing/internal/model"
)

func (h *ListingHandler) ListProvinces(c *gin.Context) {
	provinces, err := h.listingService.ListProvinces(c.Request.Context())
	if err != nil {
		log.Printf("[ERROR] failed to list provinces: %v", err)
		response.InternalServerError(c)
		return
	}

	resp := make([]ProvinceResponse, len(provinces))
	for i, p := range provinces {
		resp[i] = ProvinceResponse{
			Code:     p.Code,
			FullName: p.FullName,
		}
	}

	response.OK(c, resp, "")
}

func (h *ListingHandler) ListDistrictsByProvince(c *gin.Context) {
	code, err := strconv.ParseInt(c.Param("code"), 10, 32)
	if err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid province code format")
		return
	}
	provinceCode := int32(code)

	districts, err := h.listingService.ListDistrictsByProvince(c.Request.Context(), provinceCode)
	if err != nil {
		if errors.Is(err, model.ErrProvinceCodeNotFound) {
			response.NotFound(c, response.CodeProvinceNotFound, "Province code not found")
			return
		}

		log.Printf("[ERROR] failed to get districts by province code: %v", err)
		response.InternalServerError(c)
		return
	}

	resp := make([]DistrictResponse, len(districts))
	for i, d := range districts {
		resp[i] = DistrictResponse{
			Code:         d.Code,
			FullName:     d.FullName,
			ProvinceCode: d.ProvinceCode,
		}
	}

	response.OK(c, resp, "")
}

func (h *ListingHandler) ListWardsByDistrict(c *gin.Context) {
	code, err := strconv.ParseInt(c.Param("code"), 10, 32)
	if err != nil {
		response.BadRequest(c, response.CodeValidationFailed, "Invalid district code format")
		return
	}
	districtCode := int32(code)

	wards, err := h.listingService.ListWardsByDistrict(c.Request.Context(), districtCode)
	if err != nil {
		if errors.Is(err, model.ErrDistrictCodeNotFound) {
			response.NotFound(c, response.CodeDistrictNotFound, "District code not found")
			return
		}

		log.Printf("[ERROR] failed to get wards by district code: %v", err)
		response.InternalServerError(c)
		return
	}

	resp := make([]WardResponse, len(wards))
	for i, w := range wards {
		resp[i] = WardResponse{
			Code:         w.Code,
			FullName:     w.FullName,
			DistrictCode: w.DistrictCode,
		}
	}

	response.OK(c, resp, "")
}
