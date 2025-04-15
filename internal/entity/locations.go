package entity

type Locations struct {
	ID      int64 `gorm:"primaryKey"`
	City    string
	State   string
	Country string
	Code    string
	Lat     float64
	Lng     float64

	GeonameIDs     []LocationGeoNameIDs     `gorm:"foreignKey:LocationID"`
	AlternateNames []LocationAlternateNames `gorm:"foreignKey:LocationID"`

	VectorDistance  *float32 `gorm:"-"`
	TextMatchScore  *int64   `gorm:"-"`
	RankFusionScore *float64 `gorm:"-"`
}
