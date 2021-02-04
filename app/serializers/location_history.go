package serializers

import "time"

type LocationHistoryReq struct {
	CompanyID             uint       `json:"company_id"`
	UserID                uint       `json:"user_id"`
	CheckInTime           time.Time  `json:"check_in_time"`
	CheckOutTime          *time.Time `json:"check_out_time"`
	VisitedCompanyID      uint       `json:"visited_company_id"`
	Latitude              float64    `json:"lat"`
	Longitude             float64    `json:"lon"`
}

type LocationHistoryResp struct {
	CompanyID             uint       `json:"company_id"`
	UserID                uint       `json:"user_id"`
	CheckInTime           time.Time  `json:"check_in_time"`
	CheckOutTime          *time.Time `json:"check_out_time"`
	VisitedCompanyID      uint       `json:"visited_company_id"`
	VisitedCompanyName    string     `json:"visited_company_name"`
	VisitedCompanyAddress string     `json:"visited_company_address"`
	Latitude              float64    `json:"lat"`
	Longitude             float64    `json:"lon"`
}
