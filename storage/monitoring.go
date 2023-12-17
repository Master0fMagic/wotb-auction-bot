package storage

import (
	"context"
	"database/sql"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"sync"
)

type MonitoringStorage interface {
	Save(ctx context.Context, data dto.MonitoringData) error
	Remove(ctx context.Context, chatId int64, vehicleName string) error
	GetAll(ctx context.Context) ([]dto.MonitoringData, error)
	GetAllByVehicleAndCountGte(ctx context.Context, vehicleName string, count int) ([]dto.MonitoringData, error)
}

type RuntimeMonitoringStorage struct {
	data map[string]map[int64]dto.MonitoringData // map[vehicleName] ->| map[chatId] -> monitoring data |
	mtx  sync.RWMutex
}

func NewRuntimeMonitoringStorage() *RuntimeMonitoringStorage {
	return &RuntimeMonitoringStorage{
		data: make(map[string]map[int64]dto.MonitoringData),
		mtx:  sync.RWMutex{},
	}
}

func (s *RuntimeMonitoringStorage) Save(_ context.Context, data dto.MonitoringData) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.data[data.VehicleName] == nil {
		s.data[data.VehicleName] = map[int64]dto.MonitoringData{}
	}

	s.data[data.VehicleName][data.ChatID] = data
	return nil
}

func (s *RuntimeMonitoringStorage) Remove(_ context.Context, chatId int64, vehicleName string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.data[vehicleName], chatId)

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

func (s *RuntimeMonitoringStorage) GetAllByVehicleAndCountGte(_ context.Context, vehicleName string, count int) ([]dto.MonitoringData, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	var res []dto.MonitoringData

	for _, data := range s.data[vehicleName] {
		if data.MinimalCount >= count {
			res = append(res, data)
		}
	}

	return res, nil
}

type SQLiteMonitoringStorage struct {
	db *sql.DB
}

func NewSQLiteMonitoringStorage(dbPath string) (*SQLiteMonitoringStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &SQLiteMonitoringStorage{db: db}, nil
}

func (s *SQLiteMonitoringStorage) Save(ctx context.Context, data dto.MonitoringData) error {
	_, err := s.db.ExecContext(ctx, insertDataQuery, data.VehicleName, data.ChatID, data.MinimalCount)
	return err
}

func (s *SQLiteMonitoringStorage) Remove(ctx context.Context, chatID int64, vehicleName string) error {
	_, err := s.db.ExecContext(ctx, removeDataQuery, vehicleName, chatID)
	return err
}

func (s *SQLiteMonitoringStorage) GetAll(ctx context.Context) ([]dto.MonitoringData, error) {
	rows, err := s.db.QueryContext(ctx, getAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.MonitoringData
	for rows.Next() {
		var data dto.MonitoringData
		if err := rows.Scan(&data.VehicleName, &data.ChatID, &data.MinimalCount); err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}

func (s *SQLiteMonitoringStorage) GetAllByVehicleAndCountGte(ctx context.Context, vehicleName string, count int) ([]dto.MonitoringData, error) {
	rows, err := s.db.QueryContext(ctx, getAllByVehicleAndCountGteQuery, vehicleName, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.MonitoringData
	for rows.Next() {
		var data dto.MonitoringData
		if err := rows.Scan(&data.VehicleName, &data.ChatID, &data.MinimalCount); err != nil {
			return nil, err
		}
		result = append(result, data)
	}

	return result, nil
}
