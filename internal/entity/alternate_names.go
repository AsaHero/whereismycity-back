package entity

import "time"

type LocationAlternateNames struct {
	ID              int64 `gorm:"primaryKey"`
	LocationID      int64
	GeoNameID       int64
	AlternateNameID int64
	Type            string
	ISOLanguageCode *string
	AlternateName   string
	IsPreferred     bool
	IsShort         bool
	IsColloquial    bool
	IsHistoric      bool
	CreatedAt       time.Time
}
