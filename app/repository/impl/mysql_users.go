package impl

import (
	"clean/app/domain"
	"clean/app/repository"
	"clean/app/utils/methodutil"
	"clean/app/utils/msgutil"
	"clean/infra/conn"
	"clean/infra/errors"
	"clean/infra/logger"
	"strings"
	"time"

	stdErrors "errors"

	"gorm.io/gorm"
)

type users struct {
	*gorm.DB
}

// NewMySqlUsersRepository will create an object that represent the User.Repository implementations
func NewMySqlUsersRepository(db *gorm.DB) repository.IUsers {
	return &users{
		DB: db,
	}
}

func (r *users) Save(user *domain.User) (*domain.User, *errors.RestErr) {
	err := r.DB.Model(&domain.User{}).Create(&user).Error

	if err != nil {
		var mysqlErr conn.DbErrors

		if stdErrors.As(err, &mysqlErr.MySQLError) && mysqlErr.MySQLError.Number == errors.ErrDuplicateEntry {
			logger.Error("error occurred when create user with duplicate data", err)
			return nil, errors.NewAlreadyExistError(errors.ErrPhoneOrEmailExists)
		}
		logger.Error("error occurred when create user", err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return user, nil
}

func (r *users) GetUser(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) GetUserByID(userID uint) (*domain.User, *errors.RestErr) {
	var resp domain.User

	res := r.DB.Model(&domain.User{}).Where("id = ?", userID).First(&resp)

	if res.RowsAffected == 0 {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewNotFoundError(errors.ErrRecordNotFound)
	}

	if res.Error != nil {
		logger.Error("error occurred when getting user by user id", res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return &resp, nil
}

func (r *users) Update(user *domain.User) *errors.RestErr {
	err := r.DB.Model(&domain.User{}).Omit("password").Where("id = ?", user.ID).Updates(&user).Error

	if err != nil {
		var mysqlErr conn.DbErrors

		if stdErrors.As(err, &mysqlErr.MySQLError) && mysqlErr.MySQLError.Number == errors.ErrDuplicateEntry {
			logger.Error("error occurred when update user with duplicate data", err)
			return errors.NewAlreadyExistError(errors.ErrPhoneOrEmailExists)
		}
		logger.Error("error occurred when update user", err)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}
	return nil
}

func (r *users) UpdateUserActivation(id uint, activated bool) *errors.RestErr {

	res := r.DB.Model(&domain.User{}).Where("id = ?", id).Update("activated", !activated)

	if res.Error != nil {
		logger.Error("error occurred when updating activation of user by user id", res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) UserNameIsUnique(username string) error {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("user_name = ?", username).Find(&user)
	if res.RowsAffected != 0 {
		logger.Error("user found by this username", res.Error)
		return errors.NewError(errors.ErrUserNameNotUnique)
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by username", res.Error)
		return errors.NewError(errors.ErrSomethingWentWrong)
	}
	return nil
}

func (r *users) EmailIsUnique(email string) error {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected != 0 {
		logger.Error("user found by this email", res.Error)
		return errors.NewError(errors.ErrEmailIsUnique)
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by email", res.Error)
		return errors.NewError(errors.ErrSomethingWentWrong)
	}
	return nil
}

func (r *users) UpdatePassword(userID uint, updateValues map[string]interface{}) *errors.RestErr {
	res := r.DB.Model(&domain.User{}).Where("id = ? ", userID).Updates(&updateValues)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("updating user by user id"), res.Error)
		return errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	return nil
}

func (r *users) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	res := r.DB.Model(&domain.User{}).Where("email = ?", email).Find(&user)
	if res.RowsAffected == 0 {
		logger.Error("no user found by this email", res.Error)
		return nil, errors.NewError(errors.ErrRecordNotFound)
	}
	if res.Error != nil {
		logger.Error("error occurred when trying to get user by email", res.Error)
		return nil, errors.NewError("error occurred when trying to get user by email")
	}

	return user, nil
}

func (r *users) SetLastLoginAt(user *domain.User) error {
	lastLoginAt := time.Now().UTC()

	err := r.DB.Model(&user).Update("last_login_at", lastLoginAt).Error

	if err != nil {
		logger.Error(err.Error(), err)
		return err
	}

	return nil
}

func (r *users) HasRole(userID, roleID uint) bool {
	var count int64
	count = 0

	r.DB.Model(&domain.User{}).
		Where("id = ? AND role_id = ?", userID, roleID).
		Count(&count)

	return count > 0
}

func (r *users) ResetPassword(userID int, hashedPass []byte) error {
	err := r.DB.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("password", hashedPass).
		Error

	if err != nil {
		logger.Error("error occur when reset password", err)
		return err
	}

	return nil
}

func (r *users) GetTokenUser(id uint) (*domain.VerifyTokenResp, *errors.RestErr) {
	tempUser := &domain.TempVerifyTokenResp{}
	var vtUser domain.VerifyTokenResp

	query := r.tokenUserFetchQuery()

	res := query.Where("users.id = ?", id).Find(&tempUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("get token user"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodutil.StructToStruct(tempUser, &vtUser.BaseVerifyTokenResp)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	vtUser.Permissions = strings.Split(tempUser.Permissions, ",")

	return &vtUser, nil
}

func (r *users) GetUserWithPermissions(userID uint, withPermission bool) (*domain.UserWithPerms, *errors.RestErr) {
	var intUser domain.IntermediateUserWithPermissions
	var userWithParams domain.UserWithPerms

	sections := `
		users.*,
		roles.name role_name
	`
	if withPermission {
		sections += ",GROUP_CONCAT(DISTINCT permissions.name) AS permissions"
	}

	query := r.DB.Model(&domain.User{}).
		Select(sections).
		Joins("LEFT JOIN roles ON users.role_id = roles.id").
		Where("users.deleted_at IS NULL")

	if withPermission {
		query = query.
			Joins("JOIN role_permissions ON users.role_id = role_permissions.role_id").
			Joins("JOIN permissions ON role_permissions.permission_id = permissions.id")
	}

	query.Group("users.id")

	res := query.Where("users.id = ?", userID).Find(&intUser)

	if res.Error != nil {
		logger.Error(msgutil.EntityGenericFailedMsg("getting user with permission"), res.Error)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	err := methodutil.StructToStruct(intUser, &userWithParams.User)
	if err != nil {
		logger.Error(msgutil.EntityStructToStructFailedMsg("set intermediate user & permissions"), err)
		return nil, errors.NewInternalServerError(errors.ErrSomethingWentWrong)
	}

	userWithParams.RoleName = intUser.RoleName

	if withPermission {
		userWithParams.Permissions = strings.Split(intUser.Permissions, ",")
	}

	return &userWithParams, nil
}

func (r *users) tokenUserFetchQuery() *gorm.DB {
	selections := `
		users.id,
		users.first_name,
		users.last_name,
		users.email,
		users.phone,
		users.profile_pic,
		companies.business_id,
		businesses.name business_name,
		companies.id company_id,
		companies.name company_name,
		(
			CASE
				WHEN 1 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS admin,
		(
			CASE
				WHEN 3 IN (GROUP_CONCAT(DISTINCT users.role_id)) THEN 1 ELSE 0
			END
		) AS super_admin,
		GROUP_CONCAT(DISTINCT permissions.name) AS permissions
	`

	return r.DB.Table("users").
		Select(selections).
		Joins("LEFT JOIN companies ON users.company_id = companies.id").
		Joins("LEFT JOIN businesses ON companies.business_id = businesses.id").
		Joins("JOIN roles ON users.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.deleted_at IS NULL").
		Group("users.id")
}
