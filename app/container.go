package container

import (
	"ar5go/app/http/controllers"
	repoImpl "ar5go/app/repository/impl"
	svcImpl "ar5go/app/svc/impl"
	"ar5go/infra/conn/cache"
	"ar5go/infra/conn/db"
	"context"
)

func Init(g interface{}) {
	basectx := context.Background()
	dbc := db.Client()
	cachec := cache.Client()

	// register all repos impl, services impl, controllers
	sysRepo := repoImpl.NewSystemRepository(basectx, dbc, cachec)
	companyRepo := repoImpl.NewCompanyRepository(basectx, dbc)
	userRepo := repoImpl.NewUsersRepository(basectx, dbc)
	locationRepo := repoImpl.NewLocationRepository(basectx, dbc)
	roleRepo := repoImpl.NewRolesRepository(basectx, dbc)
	permissionRepo := repoImpl.NewPermissionsRepository(basectx, dbc)

	sysSvc := svcImpl.NewSystemService(sysRepo)
	companySvc := svcImpl.NewCompanyService(companyRepo, userRepo)
	userSvc := svcImpl.NewUsersService(basectx, userRepo)
	tokenSvc := svcImpl.NewTokenService(basectx, userRepo)
	authSvc := svcImpl.NewAuthService(basectx, userRepo, tokenSvc)
	locationSvc := svcImpl.NewLocationService(locationRepo)
	roleSvc := svcImpl.NewRolesService(roleRepo)
	permissionSvc := svcImpl.NewPermissionsService(permissionRepo)

	controllers.NewSystemController(g, sysSvc)
	controllers.NewAuthController(g, authSvc, userSvc)
	controllers.NewCompanyController(g, companySvc)
	controllers.NewUsersController(g, companySvc, userSvc, locationSvc)
	controllers.NewLocationController(g, locationSvc)
	controllers.NewRolesController(g, roleSvc)
	controllers.NewPermissionsController(g, permissionSvc)
}
