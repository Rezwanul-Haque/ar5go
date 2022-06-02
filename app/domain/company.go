package domain

import (
	"ar5go/infra/errors"
	"time"
)

type ICompany interface {
	SaveCompany(company *Company) (*Company, *errors.RestErr)
	GetCompany(companyID uint) (*Company, *errors.RestErr)
}

type Company struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Logo          string     `json:"logo"`
	Address       string     `json:"address"`
	BusinessID    uint       `json:"business_id"`
	NumOfEmployee uint       `json:"num_of_employee"`
	Email         string     `json:"email"`
	SnsLink       string     `json:"sns_link"`
	Phone         string     `json:"phone"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}
