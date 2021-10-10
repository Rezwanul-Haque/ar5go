package container

import (
	"clean/app/http/controllers"
	"clean/app/http/middlewares"
	repoImpl "clean/app/repository/impl"
	svcImpl "clean/app/svc/impl"
	"clean/infra/conn"
)

func Init(g interface{}) {
	db := conn.Db()
	mg := conn.MailGun()
	acl := middlewares.ACL

	// register all repos impl, services impl, controllers
	sysRepo := repoImpl.NewSystemRepository(db)
	uRepo := repoImpl.NewMySqlUsersRepository(db)
	roleRepo := repoImpl.NewMySqlRolesRepository(db)
	pRepo := repoImpl.NewMySqlPermissionsRepository(db)
	mRepo := repoImpl.NewMailsRepository(mg)

	sysSvc := svcImpl.NewSystemService(sysRepo)
	mailSvc := svcImpl.NewMailsService(mRepo)
	userSvc := svcImpl.NewUsersService(uRepo, mailSvc)
	tokenSvc := svcImpl.NewTokenService(uRepo)
	authSvc := svcImpl.NewAuthService(uRepo, tokenSvc)
	roleSvc := svcImpl.NewRolesService(roleRepo)
	permissionSvc := svcImpl.NewPermissionsService(pRepo)

	controllers.NewSystemController(g, sysSvc)
	controllers.NewAuthController(g, authSvc, userSvc)
	controllers.NewUsersController(g, acl, userSvc)
	controllers.NewRolesController(g, acl, roleSvc)
	controllers.NewPermissionsController(g, acl, permissionSvc)
}
