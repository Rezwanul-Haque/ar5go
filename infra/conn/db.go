package conn

import (
	"boilerplate/app/domain"
	"boilerplate/infra/config"
	"boilerplate/infra/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB

type DbErrors struct {
	*gomysql.MySQLError
}

func ConnectDb() {
	conf := config.Db()

	logger.Info("connecting to mysql at " + conf.Host + ":" + conf.Port + "...")

	logMode := gormlogger.Silent
	if conf.Debug {
		logMode = gormlogger.Info
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", conf.User, conf.Pass, conf.Host, conf.Port, conf.Schema)

	dB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormlogger.Default.LogMode(logMode),
	})

	if err != nil {
		panic(err)
	}

	sqlDb, err := dB.DB()
	if err != nil {
		panic(err)
	}

	if conf.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(conf.MaxIdleConn)
	}
	if conf.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(conf.MaxOpenConn)
	}
	if conf.MaxConnLifetime != 0 {
		sqlDb.SetConnMaxLifetime(conf.MaxConnLifetime * time.Second)
	}

	db = dB

	db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RolePermission{},
	)

	logger.Info("mysql connection successful...")
}

func Db() *gorm.DB {
	return db
}

type Seed struct {
	Name string
	Run  func(db *gorm.DB, truncate bool) error
}

func SeedAll() []Seed {
	return []Seed{
		{
			Name: "CreateRoles",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedRoles(db, "/infra/seed/roles.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "CreatePermissions",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedPermissions(db, "/infra/seed/permissions.json"); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "CreateRolePermissions",
			Run: func(db *gorm.DB, truncate bool) error {
				if err := seedRolePermissions(db, "/infra/seed/role_permissions.json"); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func seedRoles(db *gorm.DB, jsonfilPath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilPath)
	roles := []domain.Role{}

	_ = json.Unmarshal([]byte(file), &roles)

	if truncate {
		db.Exec("TRUNCATE TABLE boilerplate.role_permissions;")
		db.Exec("TRUNCATE TABLE boilerplate.permissions;")
		db.Exec("TRUNCATE TABLE boilerplate.roles;")
	}

	var count int64

	db.Model(&domain.Role{}).Count(&count)
	if count == 0 {
		db.Create(&roles)
	}

	return nil
}

func seedPermissions(db *gorm.DB, jsonfilPath string) error {
	file, _ := readSeedFile(jsonfilPath)
	perms := []domain.Permission{}

	_ = json.Unmarshal([]byte(file), &perms)

	var count int64

	db.Model(&domain.Permission{}).Count(&count)
	if count == 0 {
		db.Create(&perms)
	}

	return nil
}

func seedRolePermissions(db *gorm.DB, jsonfilPath string) error {
	file, _ := readSeedFile(jsonfilPath)
	rp := []domain.RolePermission{}

	_ = json.Unmarshal([]byte(file), &rp)

	var count int64

	db.Model(&domain.RolePermission{}).Count(&count)
	if count == 0 {
		db.Create(&rp)
	}

	return nil
}

func readSeedFile(jsonfilPath string) ([]byte, error) {
	BaseDir, _ := os.Getwd()
	seedFile := BaseDir + jsonfilPath
	if BaseDir == "/" {
		seedFile = jsonfilPath
	}
	fmt.Println("seed folder: ", seedFile)

	return ioutil.ReadFile(seedFile)
}
