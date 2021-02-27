package svc

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/serializers"
	"clean/app/svc"
	"clean/infra/errors"
	"encoding/json"
)

type history struct {
	hrepo repository.IHistory
}

func NewHistoryService(hrepo repository.IHistory) svc.IHistory {
	return &history{
		hrepo: hrepo,
	}
}

func (h *history) Create(req serializers.LocationHistoryReq) *errors.RestErr {
	var locHistory domain.LocationHistory

	jsonData, _ := json.Marshal(req)
	_ = json.Unmarshal(jsonData, &locHistory)


	err := h.hrepo.Save(&locHistory)
	if err != nil {
		return err
	}

	return nil
}
