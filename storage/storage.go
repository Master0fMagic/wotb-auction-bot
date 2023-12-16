package storage

import (
	"context"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"sync"
)

type MonitoringStorage interface {
	Save(ctx context.Context, data dto.MonitoringData) error
	Remove(ctx context.Context, chatId int, vehicleId int) error
	GetAll(ctx context.Context) ([]dto.MonitoringData, error)
	GetAllByVehicleAndCountLte(ctx context.Context, vehicleId, count int) ([]dto.MonitoringData, error)
}

type RuntimeMonitoringStorage struct {
	data map[int]map[int]dto.MonitoringData // map[vehicleId] ->| map[chatId] -> monitoring data |
	mtx  sync.RWMutex
}

func NewRuntimeMonitoringStorage() *RuntimeMonitoringStorage {
	return &RuntimeMonitoringStorage{
		data: make(map[int]map[int]dto.MonitoringData),
		mtx:  sync.RWMutex{},
	}
}

func (s *RuntimeMonitoringStorage) Save(_ context.Context, data dto.MonitoringData) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.data[data.VehicleID][data.ChatID] = data
	return nil
}

func (s *RuntimeMonitoringStorage) Remove(_ context.Context, chatId int, vehicleId int) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.data[vehicleId], chatId)

	return nil
}

func (s *RuntimeMonitoringStorage) GetAll(_ context.Context) ([]dto.MonitoringData, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	var res []dto.MonitoringData
	for _, chatMap := range s.data {
		for _, data := range chatMap {
			res = append(res, data)
		}
	}

	return res, nil
}

func (s *RuntimeMonitoringStorage) GetAllByVehicleAndCountLte(_ context.Context, vehicleId, count int) ([]dto.MonitoringData, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	var res []dto.MonitoringData

	for _, data := range s.data[vehicleId] {
		if data.MinimalCount <= count {
			res = append(res, data)
		}
	}

	return res, nil
}
