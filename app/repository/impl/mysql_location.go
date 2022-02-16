package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/infra/errors"
	"clean/infra/logger"

	"gorm.io/gorm"
)

type location struct {
	*gorm.DB
}

// NewMySqlLocationRepository will create an object that represent the Location.Repository implementations
func NewMySqlLocationRepository(db *gorm.DB) repository.ILocation {
	return &location{
		DB: db,
	}
}

func (r *location) Save(history *domain.LocationHistory) *errors.RestErr {
	res := r.DB.Model(&domain.LocationHistory{}).Create(&history)

	if res.Error != nil {
		logger.Error("error occurred when saving location history", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *location) Update(history *domain.LocationHistory) (*domain.LocationHistory, *errors.RestErr) {
	res := r.DB.Model(&domain.LocationHistory{}).
		Where("client_id = ?", history.ClientID).
		Update("check_out_time", history.CheckOutTime)

	if res.Error != nil {
		logger.Error("error occurred when saving location history", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return history, nil
}

func (r *location) GetLocationsByUserID(userID uint, pagination *serializers.Pagination) ([]*domain.IntermediateLocationHistory, *errors.RestErr) {
	var resp []*domain.IntermediateLocationHistory

	var totalRows int64 = 0
	tableName := "location_histories"
	stmt := GenerateFilteringCondition(r.DB, tableName, pagination, false)

	stmt = stmt.Model(&domain.LocationHistory{}).
		Select("location_histories.company_id, c.name as company_name, location_histories.user_id, u.user_name as name, "+
			"location_histories.id location_id, check_in_time, check_out_time, client_id, clients.name as client_name, clients.address as client_address, lat, lon").
		Joins("LEFT JOIN companies c ON c.id = location_histories.company_id").
		Joins("LEFT JOIN clients ON client_id = clients.id").
		Joins("LEFT JOIN users u ON u.id = location_histories.user_id").
		Where("location_histories.user_id = ?", userID).
		Find(&resp)

	if len(pagination.QueryString) > 0 {
		searchStmt := "clients.name LIKE ? OR clients.address LIKE ? "
		searchTerm := "%" + pagination.QueryString + "%"
		stmt.Where(searchStmt, searchTerm, searchTerm)
	}
	res := stmt.Find(&resp)
	if res.RowsAffected == 0 {
		return nil, errors.NewNotFoundError("no location histories found")
	}

	if res.Error != nil {
		logger.Error("error occurred when getting location history by userID", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	pagination.Rows = resp

	stmt = GenerateFilteringCondition(r.DB, tableName, pagination, true)
	// count all data
	errCount := r.DB.Model(&domain.LocationHistory{}).Where("location_histories.user_id = ?", userID).Count(&totalRows).Error
	if errCount != nil {
		logger.Error("error occurred when getting total location history count by userID", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	pagination.TotalRows = totalRows
	totalPages := CalculateTotalPageAndRows(pagination, totalRows)
	pagination.TotalPages = totalPages
	return resp, nil
}
