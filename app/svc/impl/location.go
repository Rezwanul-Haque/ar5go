package impl

import (
	"ar5go/app/domain"
	"ar5go/app/repository"
	"ar5go/app/serializers"
	"ar5go/app/svc"
	"ar5go/app/utils/methodsutil"
	"ar5go/app/utils/msgutil"
	"ar5go/infra/errors"
	"ar5go/infra/logger"
	"encoding/json"
)

type location struct {
	lc    logger.LogClient
	hrepo repository.ILocation
}

func NewLocationService(lc logger.LogClient, hrepo repository.ILocation) svc.ILocation {
	return &location{
		lc:    lc,
		hrepo: hrepo,
	}
}

func (h location) Create(req serializers.LocationHistoryReq) *errors.RestErr {
	var locHistory domain.LocationHistory

	jsonData, _ := json.Marshal(req)
	_ = json.Unmarshal(jsonData, &locHistory)

	err := h.hrepo.SaveLocation(&locHistory)
	if err != nil {
		return err
	}

	return nil
}

func (h location) Update(req serializers.LocationHistoryReq) (*domain.LocationHistory, *errors.RestErr) {
	var locHistory domain.LocationHistory

	jsonData, _ := json.Marshal(req)
	_ = json.Unmarshal(jsonData, &locHistory)

	resp, err := h.hrepo.UpdateLocation(&locHistory)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h location) GetLocationsByUserID(userID uint, filters *serializers.ListFilters) (*serializers.ListFilters, *errors.RestErr) {
	var resp serializers.LocationHistoryResp

	locations, err := h.hrepo.GetLocationsByUserID(userID, filters)
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
			h.lc.Error(msgutil.EntityStructToStructFailedMsg("get locations"), err)
			return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
		}
	}

	filters.Results = resp
	return filters, nil
}
