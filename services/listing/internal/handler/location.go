package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/katatrina/airbnb-clone/pkg/response"
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

	response.OK(c, resp)
}

func (h *ListingHandler) ListWardsByProvince(c *gin.Context) {
	provinceCode := c.Query("provinceCode")
	if provinceCode == "" {
		response.BadRequest(c, response.CodeValidationFailed, "provinceCode is required")
		return
	}

	wards, err := h.listingService.ListWardsByProvince(c.Request.Context(), provinceCode)
	if err != nil {
		log.Printf("[ERROR] failed to get wards by province code: %v", err)
		response.InternalServerError(c)
		return
	}

	resp := make([]WardResponse, len(wards))
	for i, w := range wards {
		resp[i] = WardResponse{
			Code:         w.Code,
			FullName:     w.FullName,
			ProvinceCode: provinceCode,
		}
	}

	response.OK(c, resp)
}
