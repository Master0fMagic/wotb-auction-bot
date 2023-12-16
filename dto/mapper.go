package dto

func MapResultToVehicleInfo(result Result) VehicleInfo {
	v := VehicleInfo{
		ID:           result.ID,
		Name:         result.Entity.UserString,
		Nation:       result.Entity.Nation,
		Type:         result.Entity.TypeSlug,
		Level:        result.Entity.RomanLevel,
		Img:          result.Entity.ImageURL,
		CurrentCount: result.CurrentCount,
		Price:        result.CurrentCount,
	}
	if result.NextPrice != nil {
		v.NextPrice = result.NextPrice.Value
	}
	return v
}
