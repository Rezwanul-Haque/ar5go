package domain

import (
	"ar5go/app/serializers"
	"ar5go/infra/errors"
	"time"
)

type ILocation interface {
	SaveLocation(company *LocationHistory) *errors.RestErr
	UpdateLocation(*LocationHistory) (*LocationHistory, *errors.RestErr)
	GetLocationsByUserID(userID uint, filters *serializers.ListFilters) ([]*IntermediateLocationHistory, *errors.RestErr)
}

type LocationHistory struct {
	ID           uint       `json:"id"`
	CompanyID    uint       `json:"company_id"`
	Company      Company    `gorm:"foreignKey:CompanyID" json:"-"`
	UserID       uint       `json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"-"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	ClientID     uint       `json:"client_id"`
	Lat          float64    `json:"lat"`
	Lon          float64    `json:"lon"`
	CreatedAt    time.Time  `json:"-"`
	UpdatedAt    time.Time  `json:"-"`
	DeletedAt    *time.Time `json:"-"`
}

type IntermediateLocationHistory struct {
	CompanyID     uint       `json:"company_id"`
	CompanyName   string     `json:"company_name"`
	UserID        uint       `json:"user_id"`
	Name          string     `json:"name"`
	LocationID    uint       `json:"location_id"`
	CheckInTime   time.Time  `json:"check_in_time"`
	CheckOutTime  *time.Time `json:"check_out_time"`
	ClientID      uint       `json:"client_id"`
	ClientName    string     `json:"client_name"`
	ClientAddress string     `json:"client_address"`
	Lat           float64    `json:"lat"`
	Lon           float64    `json:"lon"`
}

type BaseLocationHistory struct {
	CheckInTime   time.Time  `json:"check_in_time"`
	CheckOutTime  *time.Time `json:"check_out_time"`
	ClientID      uint       `json:"client_id"`
	ClientName    string     `json:"client_name"`
	ClientAddress string     `json:"client_address"`
	Lat           float64    `json:"lat"`
	Lon           float64    `json:"lon"`
}

func (*LocationHistory) TableName() string {
	return "location_histories"
}
