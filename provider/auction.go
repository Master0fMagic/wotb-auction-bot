package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"net/http"
)

type AuctionDataProvider interface {
	GetData(ctx context.Context) ([]dto.VehicleInfo, error)
}

type HttpAuctionProvider struct {
	url string
}

func NewHttpActionProvider(url string) *HttpAuctionProvider {
	return &HttpAuctionProvider{url: url}
}

func (p *HttpAuctionProvider) GetData(_ context.Context) ([]dto.VehicleInfo, error) {
	resp, err := http.Get(p.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the request was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into the Response struct
	var response dto.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	var res []dto.VehicleInfo
	for _, ent := range response.Results {
		if !ent.Available {
			continue
		}
		res = append(res, dto.MapResultToVehicleInfo(ent))
	}

	return res, nil
}
