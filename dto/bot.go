package dto

import "fmt"

type VehicleInfo struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Nation       string `json:"nation"`
	Type         string `json:"type"`
	Level        string `json:"level"`
	Img          string `json:"img"`
	CurrentCount int    `json:"current_count"`
	Price        int    `json:"price"`
	NextPrice    int    `json:"next_price,omitempty"`
}

// String returns a formatted string representation of the Vehicle
func (v *VehicleInfo) String() string {
	return fmt.Sprintf(`%s, %s %s %s
count left: %d
current price: %d; next price: %s`, v.Name, v.Level, v.Nation, v.Type, v.CurrentCount, v.Price, formatNextPrice(v.NextPrice))
}

// formatNextPrice formats the NextPrice field for display
func formatNextPrice(nextPrice int) string {
	if nextPrice == 0 {
		return "-"
	}
	return fmt.Sprintf("%d", nextPrice)
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
