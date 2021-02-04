package container

import (
	"clean/app/controllers"
	repoImpl "clean/app/repository/impl"
	svcImpl "clean/app/svc/impl"
	"clean/infrastructure/conn"
	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {
	db := conn.Db()
	redis := conn.Redis()

	// register all repos impl, services impl, controllers
	cRepo := repoImpl.NewMySqlCompanyRepository(db)
	uRepo := repoImpl.NewMySqlUsersRepository(db)
	rRepo := repoImpl.NewRedisRepository(redis)

	cSvc := svcImpl.NewCompanyService(cRepo, uRepo)
	uSvc := svcImpl.NewUsersService(uRepo)
	tSvc := svcImpl.NewTokenService(uRepo, rRepo)
	aSvc := svcImpl.NewAuthService(uRepo, rRepo, tSvc)

	controllers.NewAuthController(g, aSvc, uSvc)
	controllers.NewCompanyController(g, cSvc)
	controllers.NewUsersController(g, cSvc, uSvc)
}
