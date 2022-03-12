package db

import (
	"ar5go/app/domain"
	"ar5go/app/serializers"
	"ar5go/infra/conn/db/models"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
)

func (dc DatabaseClient) SaveLocation(history *domain.LocationHistory) *errors.RestErr {
	res := dc.DB.Model(&models.LocationHistory{}).Create(&history)

	if res.Error != nil {
		logger.Error("error occurred when saving location history", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (dc DatabaseClient) UpdateLocation(history *domain.LocationHistory) (*domain.LocationHistory, *errors.RestErr) {
	res := dc.DB.Model(&models.LocationHistory{}).
		Where("client_id = ?", history.ClientID).
		Update("check_out_time", history.CheckOutTime)

	if res.Error != nil {
		logger.Error("error occurred when saving location history", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return history, nil
}

func (dc DatabaseClient) GetLocationsByUserID(userID uint, filters *serializers.ListFilters) ([]*domain.IntermediateLocationHistory, *errors.RestErr) {
	var resp []*domain.IntermediateLocationHistory

	var totalRows int64 = 0
	tableName := "location_histories"
	stmt := applyFilteringCondition(dc.DB, tableName, filters, false)

	stmt = stmt.Model(&models.LocationHistory{}).
		Select("location_histories.company_id, c.name as company_name, location_histories.user_id, u.user_name as name, "+
			"location_histories.id location_id, check_in_time, check_out_time, client_id, clients.name as client_name, clients.address as client_address, lat, lon").
		Joins("LEFT JOIN companies c ON c.id = location_histories.company_id").
		Joins("LEFT JOIN clients ON client_id = clients.id").
		Joins("LEFT JOIN users u ON u.id = location_histories.user_id").
		Where("location_histories.user_id = ?", userID).
		Find(&resp)

	if len(filters.QueryString) > 0 {
		searchStmt := "clients.name LIKE ? OR clients.address LIKE ? "
		searchTerm := "%" + filters.QueryString + "%"
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

	filters.Rows = resp

	stmt = applyFilteringCondition(dc.DB, tableName, filters, true)
	// count all data
	errCount := dc.DB.Model(&models.LocationHistory{}).Where("location_histories.user_id = ?", userID).Count(&totalRows).Error
	if errCount != nil {
		logger.Error("error occurred when getting total location history count by userID", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	filters.TotalRows = totalRows
	filters.CalculateTotalPageAndRows(totalRows)

	return resp, nil
}
