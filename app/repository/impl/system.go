package impl

import (
	"boilerplate/app/repository"

	"gorm.io/gorm"
)

type system struct {
	*gorm.DB
}

// NewSystemRepository will create an object that represent the System.Repository implementations
func NewSystemRepository(db *gorm.DB) repository.ISystem {
	return &system{
		DB: db,
	}
}

func (sys *system) DBCheck() (bool, error) {
	dB, _ := sys.DB.DB()
	if err := dB.Ping(); err != nil {
		return false, err
	}

	return true, nil
}
