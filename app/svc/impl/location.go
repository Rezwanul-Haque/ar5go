package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/app/utils/methodsutil"
	"clean/app/utils/msgutil"
	"clean/infra/errors"
	"clean/infra/logger"
	"encoding/json"
)

type location struct {
	hrepo repository.ILocation
}

func NewLocationService(hrepo repository.ILocation) svc.ILocation {
	return &location{
		hrepo: hrepo,
	}
}

func (h *location) Create(req serializers.LocationHistoryReq) *errors.RestErr {
	var locHistory domain.LocationHistory

	jsonData, _ := json.Marshal(req)
	_ = json.Unmarshal(jsonData, &locHistory)

	err := h.hrepo.Save(&locHistory)
	if err != nil {
		return err
	}

	return nil
}

func (h *location) Update(req serializers.LocationHistoryReq) (*domain.LocationHistory, *errors.RestErr) {
	var locHistory domain.LocationHistory

	jsonData, _ := json.Marshal(req)
	_ = json.Unmarshal(jsonData, &locHistory)

	resp, err := h.hrepo.Update(&locHistory)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *location) GetLocationsByUserID(userID uint, pagination *serializers.Pagination) (*serializers.Pagination, *errors.RestErr) {
	var resp serializers.LocationHistoryResp

	locations, err := h.hrepo.GetLocationsByUserID(userID, pagination)
	if err != nil {
		return nil, err
	}

	resp.CompanyID = locations[0].CompanyID
	resp.CompanyName = locations[0].CompanyName
	resp.UserID = locations[0].UserID
	resp.UserName = locations[0].Name

	if locations[0].LocationID > 0 {
		err := methodsutil.StructToStruct(locations, &resp.Locations)
		if err != nil {
			logger.Error(msgutil.EntityStructToStructFailedMsg("get locations"), err)
			return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		}
	}

	pagination.Rows = resp
	return pagination, nil
}
