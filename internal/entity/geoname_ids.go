package entity

import "time"

type LocationGeoNameIDs struct {
	LocationID uint `gorm:"primaryKey"`
	GeoNameID  uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}
