package seeder

import (
	"ar5go/infra/conn/db"
	"ar5go/infra/conn/db/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const SeedFilesPath = "/_fixture/seed"

type Seed struct {
	Name string
	Run  func(dbc db.DatabaseClient, truncate bool) error
}

func SeedAll() []Seed {
	return []Seed{
		{
			Name: "CreateRoles",
			Run: func(dbc db.DatabaseClient, truncate bool) error {
				if err := seedRoles(dbc, SeedFilesPath+"/roles.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "CreatePermissions",
			Run: func(dbc db.DatabaseClient, truncate bool) error {
				if err := seedPermissions(dbc, SeedFilesPath+"/permissions.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name: "CreateRolePermissions",
			Run: func(dbc db.DatabaseClient, truncate bool) error {
				if err := seedRolePermissions(dbc, SeedFilesPath+"/role_permissions.json", truncate); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

func seedRoles(dbc db.DatabaseClient, jsonfilePath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilePath)
	var roles []models.Role

	_ = json.Unmarshal(file, &roles)

	if truncate {
		dbc.DB.Exec("TRUNCATE TABLE ar5go_db.role_permissions;")
		dbc.DB.Exec("TRUNCATE TABLE ar5go_db.permissions;")
		dbc.DB.Exec("TRUNCATE TABLE ar5go_db.roles;")
	}

	var count int64

	dbc.DB.Model(&models.Role{}).Count(&count)
	if count == 0 {
		dbc.DB.Create(&roles)
	}

	return nil
}

func seedPermissions(dbc db.DatabaseClient, jsonfilePath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilePath)
	var perms []models.Permission

	_ = json.Unmarshal(file, &perms)

	var count int64

	dbc.DB.Model(&models.Permission{}).Count(&count)
	if count == 0 {
		dbc.DB.Create(&perms)
	}

	return nil
}

func seedRolePermissions(dbc db.DatabaseClient, jsonfilePath string, truncate bool) error {
	file, _ := readSeedFile(jsonfilePath)
	var rp []models.RolePermission

	_ = json.Unmarshal(file, &rp)

	var count int64

	dbc.DB.Model(&models.RolePermission{}).Count(&count)
	if count == 0 {
		dbc.DB.Create(&rp)
	}

	return nil
}

func readSeedFile(jsonfilePath string) ([]byte, error) {
	BaseDir, _ := os.Getwd()
	seedFile := BaseDir + jsonfilePath
	if BaseDir == "/" {
		seedFile = jsonfilePath
	}
	fmt.Println("seed folder: ", seedFile)

	return ioutil.ReadFile(seedFile)
}
