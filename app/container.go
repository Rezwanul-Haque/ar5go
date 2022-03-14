package container

import (
	"ar5go/app/http/controllers"
	repoImpl "ar5go/app/repository/impl"
	svcImpl "ar5go/app/svc/impl"
	"ar5go/infra/conn/cache"
	"ar5go/infra/conn/db"
	"ar5go/infra/logger"
	"context"
)

func Init(g interface{}, lc logger.LogClient) {
	basectx := context.Background()
	dbc := db.Client()
	cachec := cache.Client()

	// register all repos impl, services impl, controllers
	sysRepo := repoImpl.NewSystemRepository(basectx, lc, dbc, cachec)
	companyRepo := repoImpl.NewCompanyRepository(basectx, lc, dbc)
	userRepo := repoImpl.NewUsersRepository(basectx, lc, dbc)
	locationRepo := repoImpl.NewLocationRepository(basectx, lc, dbc)
	roleRepo := repoImpl.NewRolesRepository(basectx, lc, dbc)
	permissionRepo := repoImpl.NewPermissionsRepository(basectx, lc, dbc)

	sysSvc := svcImpl.NewSystemService(sysRepo)
	companySvc := svcImpl.NewCompanyService(lc, companyRepo, userRepo)
	userSvc := svcImpl.NewUsersService(basectx, lc, userRepo)
	tokenSvc := svcImpl.NewTokenService(basectx, lc, userRepo)
	authSvc := svcImpl.NewAuthService(basectx, lc, userRepo, tokenSvc)
	locationSvc := svcImpl.NewLocationService(lc, locationRepo)
	roleSvc := svcImpl.NewRolesService(lc, roleRepo)
	permissionSvc := svcImpl.NewPermissionsService(lc, permissionRepo)

	controllers.NewSystemController(g, lc, sysSvc)
	controllers.NewAuthController(g, lc, authSvc, userSvc)
	controllers.NewCompanyController(g, lc, companySvc)
	controllers.NewUsersController(g, lc, companySvc, userSvc, locationSvc)
	controllers.NewLocationController(g, lc, locationSvc)
	controllers.NewRolesController(g, lc, roleSvc)
	controllers.NewPermissionsController(g, lc, permissionSvc)
}
