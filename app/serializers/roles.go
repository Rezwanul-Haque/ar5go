package serializers

import v "github.com/go-ozzo/ozzo-validation/v4"

type RoleReq struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type RoleResp struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

type RolePermissionsReq struct {
	RoleID      int   `json:"role_id,omitempty"`
	Permissions []int `json:"permissions"`
}

func (r RoleReq) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Name, v.Required, v.Length(5, 0)),
		v.Field(&r.DisplayName, v.Required, v.Length(5, 0)),
	)
}

func (rp RolePermissionsReq) Validate() error {
	return v.ValidateStruct(&rp,
		v.Field(&rp.RoleID),
		v.Field(&rp.Permissions, v.Required, v.Length(1, 0)),
	)
}
