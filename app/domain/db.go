package domain

type IDb interface {
	ICompany
	IUsers
	ILocation
	IPermissions
	IRoles
}
