package storage

import (
	"context"
	"errors"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"sync"
)

type AddMonitoringFlowStorage interface {
	Set(ctx context.Context, data dto.AddMonitoringStep) error
	Get(ctx context.Context, chatID int64) (*dto.AddMonitoringStep, error)
	Remove(ctx context.Context, chatID int64) error
}

type RuntimeAddMonitoringFlowStorage struct {
	mtx  sync.RWMutex
	data map[int64]dto.AddMonitoringStep // chatID -> data
}

func NewRuntimeAddMonitoringFlowStorage() *RuntimeAddMonitoringFlowStorage {
	return &RuntimeAddMonitoringFlowStorage{
		data: make(map[int64]dto.AddMonitoringStep),
		mtx:  sync.RWMutex{},
	}
}

func (s *RuntimeAddMonitoringFlowStorage) Set(_ context.Context, data dto.AddMonitoringStep) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.data[data.Data.ChatID] = data
	return nil
}

func (s *RuntimeAddMonitoringFlowStorage) Get(_ context.Context, chatID int64) (*dto.AddMonitoringStep, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	data, ok := s.data[chatID]
	if !ok {
		return nil, errors.New("step not found")
	}
	return &data, nil
}

func (s *RuntimeAddMonitoringFlowStorage) Remove(_ context.Context, chatID int64) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.data, chatID)
	return nil
}
