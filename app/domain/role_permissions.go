package domain

type RolePermission struct {
	ID           uint `gorm:"primaryKey;autoIncrement:true" json:"id"`
	RoleID       uint `gorm:"primaryKey;autoIncrement:false" json:"role_id"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false" json:"permission_id"`
	// putting this relations make the seeder unusable
	// Role         Role       `gorm:"foreignKey:RoleID" json:"-"`
	// Permission   Permission `gorm:"foreignKey:PermissionID" json:"-"`
}

type RolePermissions struct {
	RoleID      uint   `json:"role_id"`
	Permissions []uint `json:"permissions"`
}
