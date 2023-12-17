package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"net/http"
	"sync"
	"time"
)

type AuctionDataProvider interface {
	GetData(ctx context.Context, skipSold bool) ([]dto.VehicleInfo, error)
}

type HttpAuctionProvider struct {
	url string
}

func NewHttpActionProvider(url string) *HttpAuctionProvider {
	return &HttpAuctionProvider{url: url}
}

func (p *HttpAuctionProvider) GetData(_ context.Context, skipSold bool) ([]dto.VehicleInfo, error) {
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
		if skipSold && ent.CurrentCount == 0 {
			continue
		}
		res = append(res, dto.MapResultToVehicleInfo(ent))
	}

	return res, nil
}

type CachedAuctionDataProvider struct {
	data []dto.VehicleInfo
	mtx  sync.Mutex

	cacheLifetime time.Duration
	dataProvider  AuctionDataProvider
}

func NewCachedAuctionDataProvider(dataProvider AuctionDataProvider, cacheLifetime time.Duration) *CachedAuctionDataProvider {
	return &CachedAuctionDataProvider{
		data:          make([]dto.VehicleInfo, 0),
		mtx:           sync.Mutex{},
		cacheLifetime: cacheLifetime,
		dataProvider:  dataProvider,
	}
}

func (p *CachedAuctionDataProvider) Run(ctx context.Context) error {
	ticker := time.NewTicker(p.cacheLifetime)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			p.mtx.Lock()
			p.data = make([]dto.VehicleInfo, 0)
			p.mtx.Unlock()
		}
	}
}

func (p *CachedAuctionDataProvider) GetData(ctx context.Context, skipSold bool) ([]dto.VehicleInfo, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if len(p.data) == 0 {
		data, err := p.dataProvider.GetData(ctx, false)
		if err != nil {
			return nil, err
		}
		p.data = data
	}

	var res []dto.VehicleInfo
	for _, v := range p.data {
		if skipSold && v.CurrentCount == 0 {
			continue
		}
		res = append(res, v)
	}

	return res, nil
}
