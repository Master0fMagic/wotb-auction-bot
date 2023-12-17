package dto

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func MapResultToVehicleInfo(result Result) VehicleInfo {
	v := VehicleInfo{
		ID:            result.ID,
		Name:          result.Entity.UserString,
		Nation:        result.Entity.Nation,
		Type:          result.Entity.TypeSlug,
		Level:         result.Entity.RomanLevel,
		Img:           result.Entity.ImageURL,
		InitialCount:  result.InitialCount,
		CurrentCount:  result.CurrentCount,
		Price:         result.Price.Value,
		AvailableTill: result.AvailableBefore,
		Available:     result.Available,
	}
	if result.NextPrice != nil {
		v.NextPrice = result.NextPrice.Value
	}

	if t, err := time.Parse("2006-01-02T15:04:05", result.AvailableBefore); err != nil {
		log.WithError(err).Error("error parsing date. Skip formatting")
	} else {
		v.AvailableTill = t.Format(time.RFC1123)
	}

	return v
}
