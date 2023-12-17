package dto

import (
	"fmt"
)

type VehicleInfo struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Nation        string `json:"nation"`
	Type          string `json:"type"`
	Level         string `json:"level"`
	Img           string `json:"img"`
	InitialCount  int    `json:"initialCount"`
	CurrentCount  int    `json:"current_count"`
	Price         int    `json:"price"`
	NextPrice     int    `json:"next_price,omitempty"`
	Available     bool   `json:"available"`
	AvailableTill string `json:"available_till"`
}

// String returns a formatted string representation of the Vehicle
func (v *VehicleInfo) String() string {
	return fmt.Sprintf(`%s, %s %s %s
status: %s
initial count: %d, count left: %d
current price: %d; next price: %s
available till: %s`, v.Name, v.Level, v.Nation, v.Type,
		v.getStatus(), v.InitialCount, v.CurrentCount, v.Price, formatNextPrice(v.NextPrice), v.AvailableTill)
}

// formatNextPrice formats the NextPrice field for display
func formatNextPrice(nextPrice int) string {
	if nextPrice == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", nextPrice)
}

func (v *VehicleInfo) getStatus() string {
	if v.CurrentCount == 0 {
		return "💲 sold 💲"
	}
	if !v.Available {
		return "❌ sale is over ❌"
	}
	return "✅ available ✅"
}

type MonitoringData struct {
	VehicleName  string `json:"vehicle_name"`
	ChatID       int64  `json:"chat_id"`
	MinimalCount int    `json:"minimal_count"`
}

type MonitoringStep int

const (
	StepSelectVehicle MonitoringStep = iota
	StepEnterMinimalCount
)

type AddMonitoringStep struct {
	Data MonitoringData
	Step MonitoringStep
}
