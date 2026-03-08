package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/katatrina/airbnb-clone/services/booking/internal/model"
	"github.com/katatrina/airbnb-clone/services/booking/internal/service"
)

type ListingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewListingClient(baseURL string) *ListingClient {
	return &ListingClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// listingAPIResponse maps to Listing Service's actual JSON response structure.
// This struct is PRIVATE — only used inside this package for JSON parsing.
type listingAPIResponse struct {
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Data    *listingAPIData `json:"data"`
}

type listingAPIData struct {
	ID            string `json:"id"`
	HostID        string `json:"hostId"`
	PricePerNight int64  `json:"pricePerNight"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
}

func (c *ListingClient) GetActiveListingByID(
	ctx context.Context,
	id string,
) (*service.ListingInfo, error) {
	url := fmt.Sprintf("%s/api/v1/listings/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, model.ErrListingServiceUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, model.ErrListingNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, model.ErrListingServiceUnavailable
	}

	var apiResp listingAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode listing response: %w", err)
	}

	if apiResp.Data == nil {
		return nil, model.ErrListingNotFound
	}

	return &service.ListingInfo{
		ID:            apiResp.Data.ID,
		HostID:        apiResp.Data.HostID,
		PricePerNight: apiResp.Data.PricePerNight,
		Currency:      apiResp.Data.Currency,
		Status:        apiResp.Data.Status,
	}, nil
}
