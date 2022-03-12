package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/infra/conn/db"
	"ar5go/infra/errors"
	"context"
)

type company struct {
	ctx context.Context
	DB  db.DatabaseClient
}

// NewCompanyRepository will create an object that represent the Company.Repository implementations
func NewCompanyRepository(ctx context.Context, dbc db.DatabaseClient) repository.ICompany {
	return &company{
		ctx: ctx,
		DB:  dbc,
	}
}

func (r *company) SaveCompany(com *domain.Company) (*domain.Company, *errors.RestErr) {
	return r.DB.SaveCompany(com)
}

func (r *company) GetCompany(companyID uint) (*domain.Company, *errors.RestErr) {
	return r.DB.GetCompany(companyID)
}
