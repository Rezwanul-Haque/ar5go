package serializers

import (
	"strings"
)

type CompanyReq struct {
	Name          string  `json:"name"`
	Logo          string  `json:"logo"`
	Address       string  `json:"address"`
	BusinessID    uint    `json:"business_id"`
	UserName      string  `json:"user_name"`
	Password      *string `json:"password,omitempty"`
	NumOfEmployee uint    `json:"num_of_employee"`
	Email         string  `json:"email"`
	SnsLink       string  `json:"sns_link"`
	Phone         string  `json:"phone"`
}

type CompanyResp struct {
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	Address       string `json:"address"`
	BusinessID    uint   `json:"business_id"`
	NumOfEmployee uint   `json:"num_of_employee"`
	Email         string `json:"email"`
	SnsLink       string `json:"sns_link"`
	Phone         string `json:"phone"`
	Admin         User   `json:"admin"`
}

func (r *CompanyReq) TrimRequestBody() {
	r.Email = strings.TrimSpace(r.Email)
	r.Name = strings.TrimSpace(r.Name)
}
