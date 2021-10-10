package consts

const (
	RoleIDAdmin      uint = 1
	RoleIDSales      uint = 2
	RoleIDSuperAdmin uint = 3

	RoleAdmin      string = "admin"
	RoleSales      string = "sales"
	RoleSuperAdmin string = "super_admin"

	AccessTokenType  = "access"
	RefreshTokenType = "refresh"

	PermissionUserCreate            = "create.user"
	PermissionClientCreate          = "create.client"
	PermissionLocationCreate        = "create.location"
	PermissionProductCreate         = "create.product"
	PermissionCompanyInterestCreate = "create.interest"

	PermissionRoleCrud       = "crud.role"
	PermissionRoleFetchAll   = "fetch.role.all"
	PermissionPermissionCrud = "crud.permission"

	PermissionUserUpdate            = "update.user"
	PermissionClientUpdate          = "update.client"
	PermissionCompanyUpdate         = "update.company"
	PermissionLocationUpdate        = "update.location"
	PermissionProductUpdate         = "update.product"
	PermissionCompanyInterestUpdate = "update.interest"

	PermissionUserDelete            = "delete.user"
	PermissionClientDelete          = "delete.client"
	PermissionLocationDelete        = "delete.location"
	PermissionProductDelete         = "delete.product"
	PermissionCompanyInterestDelete = "delete.interest"

	PermissionUserFetch            = "fetch.user"
	PermissionClientFetch          = "fetch.client"
	PermissionLocationFetch        = "fetch.location"
	PermissionUserLocationFetch    = "fetch.user.location"
	PermissionProductFetch         = "fetch.product"
	PermissionCompanyInterestFetch = "fetch.interest"
	PermissionImagesFetch          = "fetch.images"

	PermissionUserFetchAll            = "fetch.user.all"
	PermissionClientFetchAll          = "fetch.client.all"
	PermissionLocationFetchAll        = "fetch.location.all"
	PermissionProductFetchAll         = "fetch.product.all"
	PermissionCompanyInterestFetchAll = "fetch.interest.all"

	PermissionImagesUpload = "upload.images"

	PermissionUserUpdateManually = "update.user.manually"
)
