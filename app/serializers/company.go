package serializers

import (
	"clean/app/domain"
	"strings"
)

type CompanyPayload struct {
	Name          string  `json:"name"`
	Logo          string  `json:"logo"`
	Address       string  `json:"address"`
	BusinessID    uint    `json:"business_id"`
	Password      *string `json:"password,omitempty"`
	NumOfEmployee uint    `json:"num_of_employee"`
	Email         string  `json:"email"`
	SnsLink       string  `json:"sns_link"`
	Phone         string  `json:"phone"`
}

type CompanyResponse struct {
	Name          string      `json:"name"`
	Logo          string      `json:"logo"`
	Address       string      `json:"address"`
	BusinessID    uint        `json:"business_id"`
	NumOfEmployee uint        `json:"num_of_employee"`
	Email         string      `json:"email"`
	SnsLink       string      `json:"sns_link"`
	Phone         string      `json:"phone"`
	Admin         domain.User `json:"admin"`
}

func (r *CompanyPayload) TrimRequestBody() {
	r.Email = strings.TrimSpace(r.Email)
	r.Name = strings.TrimSpace(r.Name)
}
