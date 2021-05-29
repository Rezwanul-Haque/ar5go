package serializers

import "time"

type LocationHistoryReq struct {
	CompanyID    uint       `json:"company_id"`
	UserID       uint       `json:"user_id"`
	CheckInTime  time.Time  `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
	ClientID     uint       `json:"client_id"`
	Lat          float64    `json:"lat"`
	Lon          float64    `json:"lon"`
}

type BaseLocationHistory struct {
	LocationID    uint       `json:"location_id"`
	CheckInTime   time.Time  `json:"check_in_time"`
	CheckOutTime  *time.Time `json:"check_out_time"`
	ClientID      uint       `json:"client_id"`
	ClientName    string     `json:"client_name"`
	ClientAddress string     `json:"client_address"`
	Lat           float64    `json:"lat"`
	Lon           float64    `json:"lon"`
}

type LocationHistoryResp struct {
	CompanyID   uint                   `json:"company_id"`
	CompanyName string                 `json:"company_name"`
	UserID      uint                   `json:"user_id"`
	UserName    string                 `json:"user_name"`
	Locations   []*BaseLocationHistory `json:"locations"`
}
