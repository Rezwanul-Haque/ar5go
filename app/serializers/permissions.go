package serializers

import v "github.com/go-ozzo/ozzo-validation/v4"

type PermissionReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionResp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *PermissionReq) Validate() error {
	return v.ValidateStruct(p,
		v.Field(&p.Name, v.Required, v.Length(5, 0)),
	)
}
